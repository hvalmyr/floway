<script setup lang="ts">
import { useFocusTrap } from "@vueuse/integrations/useFocusTrap";
import { ChevronLeft, ChevronRight, X } from "lucide-vue-next";
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from "vue";
import type { GalleryPhoto } from "~/types/api";

/**
 * Autoplaying, infinitely-looping row of vertical photos (admin-managed via
 * /admin/gallery-photos), full site width — cards are 2/5 of the viewport
 * tall (`h-[40vh]`, see the slide's height class below) at every
 * breakpoint, with however many fit at that height showing side by side;
 * width per card follows from its `aspect-[3/4]` box, not a fixed
 * items-per-row count. Sliding moves one photo at a time, not one full
 * screen at a time.
 *
 * Looping uses the classic "triple the list" trick: the track is
 * [...photos, ...photos, ...photos] and starts on the middle copy. Stepping
 * ±1 through a repeated sequence looks identical to stepping through the
 * real one even while several clone/real slides are visible at once mid-
 * transition, so — unlike a single boundary clone — this stays seamless
 * regardless of how many items a given breakpoint shows at a time. Once the
 * index drifts into the leading or trailing copy, the track jumps by
 * ±photos.length once the slide transition has had time to finish (a timer
 * matched to the CSS transition duration, not `transitionend`, since that
 * event never fires under `prefers-reduced-motion`).
 *
 * The step size (one slide's width, including the gap) is measured off the
 * real DOM rather than computed analytically — every card renders at the
 * same width (same height + same `aspect-[3/4]`, cropped via `object-cover`
 * regardless of a given photo's real aspect ratio), so reading it off the
 * first slide is exact for all of them, and it stays correct if the height
 * class ever grows responsive breakpoints again without duplicating them
 * into JS.
 * The track offset (`offsetPx`) centers the active slide in the visible
 * viewport rather than pinning it to the left edge, so neighbors crop
 * evenly on both sides instead of only the trailing one being cut off.
 *
 * Also draggable — press and hold (mouse) or touch-and-drag (finger) to
 * scroll it directly, same as a native horizontal strip, rather than only
 * a single step per swipe. See `startDrag`/`moveDrag`/`endDrag` for the
 * shared core both input types feed into. The track is `touch-pan-y`
 * (not the CSS default `auto`) so the browser's own gesture recognizer
 * commits to vertical panning immediately and leaves horizontal motion to
 * `onTouchMove`'s axis lock — without it, mobile Safari/Chrome fight the JS
 * drag for the first several frames of every gesture (visible stutter) and,
 * worse, can read a drag starting near the screen edge as its own
 * back/forward swipe gesture instead of scrolling the carousel.
 *
 * Autoplay pauses for as long as the user is hovering, has keyboard focus
 * inside, is dragging, or has a photo open fullscreen, and resumes the
 * moment they're not — never a fixed cooldown. The prev/next *buttons*
 * (not autoplay, not dragging) do get a short fixed cooldown of their own
 * (`NAV_COOLDOWN_MS`) purely to stop mashing them from outrunning the
 * slide transition — see the comment by that constant.
 *
 * Clicking any photo opens it fullscreen in a lightbox (own prev/next +
 * Escape/backdrop-click to close, focus-trapped while open), filling most
 * of the screen (90vh tall). Clicking the photo again zooms it to its real
 * resolution with scrollbars to pan around — see the template comment
 * above the lightbox `<img>` for why the zoom wrapper has to stay
 * unpositioned.
 *
 * @example
 * <PhotoCarousel :photos="galleryPhotos" />
 */
const props = defineProps<{ photos: GalleryPhoto[] }>();

const AUTOPLAY_MS = 4000;
// Below this, a press-and-release (mouse or touch) counts as a tap/click
// rather than a drag — see openLightbox's `didDrag` check.
const DRAG_CLICK_THRESHOLD_PX = 6;
// Cards in the row are tall (up to 66vh, see the slide's height classes
// below) but narrow (aspect-[3/4]) — no reason to make the browser pull
// down the same multi-hundred-KB full-size file (the one the lightbox/
// zoom need) when a lighter rendition looks identical at card size. Fed
// to NuxtImg's `sizes` (not `width`) — `width` alone makes it emit a bare
// 1x/2x density srcset with no `sizes` attribute, so the browser picks by
// screen DPR alone with no idea the image only renders this small; a
// real `sizes` value lets it size correctly against the DPR it actually
// has (confirmed live: bare `width` was shipping 1800px files to a ~300px
// slide on a real mobile audit).
const THUMBNAIL_WIDTH = 400;
// Matches the track's `duration-500` class, plus a small buffer so the
// snap-back never fires before the (possibly reduced-motion-skipped) CSS
// transition has actually finished.
const TRANSITION_MS = 500;

const loops = computed(() => props.photos.length > 1);
const slides = computed(() =>
  loops.value ? [...props.photos, ...props.photos, ...props.photos] : props.photos,
);

const index = ref(props.photos.length);
const jumping = ref(false);
const paused = ref(false);
const trackRef = ref<HTMLElement | null>(null);
const viewportRef = ref<HTMLElement | null>(null);
const stepPx = ref(0);
const viewportWidth = ref(0);

// Centers the active slide instead of pinning it to the left edge — with
// several cards visible at once, left-pinning left the active one flush
// against the left edge and only the slide at the *right* edge partially
// cropped, an asymmetry with no clear "this one's current" focal point.
// Centering crops evenly on both sides and puts the active photo in the
// middle, matching how center-mode carousels usually read.
const offsetPx = computed(
  () => index.value * stepPx.value + stepPx.value / 2 - viewportWidth.value / 2,
);

// Press-and-hold-to-drag (mouse or touch) — the track follows the
// pointer 1:1 while `dragging`, rather than only reacting once the
// gesture ends. `dragDeltaPx` is subtracted so that dragging left
// (negative delta) increases the offset, moving the track left to reveal
// later slides, matching the direction the pointer actually moved.
const dragging = ref(false);
const dragDeltaPx = ref(0);
const displayOffsetPx = computed(() =>
  dragging.value ? offsetPx.value - dragDeltaPx.value : offsetPx.value,
);

const activeDot = computed(() =>
  loops.value
    ? (((index.value - props.photos.length) % props.photos.length) + props.photos.length) %
      props.photos.length
    : 0,
);

let timer: ReturnType<typeof setInterval> | null = null;
let snapTimer: ReturnType<typeof setTimeout> | null = null;
let resetTimer: ReturnType<typeof setTimeout> | null = null;
let resizeObserver: ResizeObserver | null = null;

function measureStep() {
  if (viewportRef.value) viewportWidth.value = viewportRef.value.clientWidth;
  const track = trackRef.value;
  const first = track?.children[0] as HTMLElement | undefined;
  if (!track || !first) return;
  const gap = parseFloat(getComputedStyle(track).columnGap || "0");
  stepPx.value = first.getBoundingClientRect().width + gap;
}

function next() {
  index.value++;
}
function prev() {
  index.value--;
}
function goTo(i: number) {
  index.value = loops.value ? i + props.photos.length : i;
}

// Arrow-button-only throttle — autoplay and swipe call next()/prev()
// directly and stay unaffected. Must be at least TRANSITION_MS: a shorter
// cooldown let a click land while the previous slide transition was still
// interpolating, interrupting it mid-flight and reading as two slides
// jumping at once instead of one clean step — exactly what let fast
// clicking outrun the row into not-yet-visible cards. The +50ms buffer
// covers the timer firing a frame or two late.
const NAV_COOLDOWN_MS = TRANSITION_MS + 50;
const navCooldown = ref(false);
function onArrowClick(action: () => void) {
  if (navCooldown.value) return;
  action();
  navCooldown.value = true;
  setTimeout(() => {
    navCooldown.value = false;
  }, NAV_COOLDOWN_MS);
}

function thumbUrl(photo: GalleryPhoto): string {
  return resolveOptimizedMediaUrl(photo.image);
}

// See the component doc comment: the "danger zone" is the whole leading/
// trailing copy now, not a single boundary index, so this checks a range
// rather than an exact edge value.
watch(index, (value) => {
  if (snapTimer !== null) {
    clearTimeout(snapTimer);
    snapTimer = null;
  }
  if (!loops.value) return;
  const length = props.photos.length;
  if (value >= length && value < length * 2) return;

  snapTimer = setTimeout(() => {
    snapTimer = null;
    jumping.value = true;
    index.value = value < length ? value + length : value - length;
    if (resetTimer !== null) clearTimeout(resetTimer);
    // Re-enable the transition on a short timer, not requestAnimationFrame —
    // rAF can be throttled or never fire at all in some contexts (e.g. a
    // backgrounded tab), which would leave the track permanently
    // transition-less. A plain timer only needs the JS event loop.
    resetTimer = setTimeout(() => {
      resetTimer = null;
      jumping.value = false;
    }, 20);
  }, TRANSITION_MS);
});

watch(
  () => props.photos.length,
  () => nextTick(measureStep),
);

function startAutoplay() {
  stopAutoplay();
  if (!loops.value) return;
  timer = setInterval(() => {
    if (!paused.value) next();
  }, AUTOPLAY_MS);
}
function stopAutoplay() {
  if (timer !== null) {
    clearInterval(timer);
    timer = null;
  }
}

function pause() {
  paused.value = true;
}
function resume() {
  paused.value = false;
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === "ArrowLeft") {
    event.preventDefault();
    prev();
  } else if (event.key === "ArrowRight") {
    event.preventDefault();
    next();
  }
}

// Shared press-and-drag core for both mouse and touch — see `dragging`/
// `dragDeltaPx` above for how this feeds the live-follow offset. Steps are
// rounded from however far the drag actually traveled (in stepPx units)
// rather than a fixed one-slide-per-gesture jump, so a long drag can
// cross several slides at once, same as a native scrollable strip.
let dragStartX = 0;
let didDrag = false;

function startDrag(clientX: number) {
  dragging.value = true;
  dragStartX = clientX;
  dragDeltaPx.value = 0;
  didDrag = false;
  pause();
}
function moveDrag(clientX: number) {
  if (!dragging.value) return;
  dragDeltaPx.value = clientX - dragStartX;
  if (Math.abs(dragDeltaPx.value) > DRAG_CLICK_THRESHOLD_PX) didDrag = true;
}
function endDrag() {
  if (!dragging.value) return;
  dragging.value = false;
  const steps = stepPx.value > 0 ? Math.round(-dragDeltaPx.value / stepPx.value) : 0;
  dragDeltaPx.value = 0;
  if (steps !== 0) index.value += steps;
  resume();
}

function onMouseDown(event: MouseEvent) {
  // Stops the browser's native "drag this image" ghost/cursor, which would
  // otherwise fight with our own drag and briefly show a floating copy of
  // the photo under the pointer.
  event.preventDefault();
  startDrag(event.clientX);
  window.addEventListener("mousemove", onMouseMove);
  window.addEventListener("mouseup", onMouseUp);
}
function onMouseMove(event: MouseEvent) {
  moveDrag(event.clientX);
}
function onMouseUp() {
  endDrag();
  window.removeEventListener("mousemove", onMouseMove);
  window.removeEventListener("mouseup", onMouseUp);
}

// Touch adds one thing mouse doesn't need: axis locking. A finger coming
// down on the carousel might mean "drag the carousel" or "scroll the
// page" — that's only knowable once it's moved a bit, so the first few
// pixels of movement decide which axis owns the gesture. Only once locked
// horizontal do we `preventDefault()` (which is also *why* touchmove can't
// be `.passive` like touchstart/touchend still are); locked vertical hands
// the gesture back to the page by ending our drag early.
let touchStartY = 0;
let touchAxis: "x" | "y" | null = null;

function onTouchStart(event: TouchEvent) {
  const touch = event.touches[0]!;
  touchStartY = touch.clientY;
  touchAxis = null;
  startDrag(touch.clientX);
}
function onTouchMove(event: TouchEvent) {
  if (!dragging.value) return;
  const touch = event.touches[0]!;
  if (touchAxis === null) {
    const dx = touch.clientX - dragStartX;
    const dy = touch.clientY - touchStartY;
    if (Math.abs(dx) < DRAG_CLICK_THRESHOLD_PX && Math.abs(dy) < DRAG_CLICK_THRESHOLD_PX) return;
    touchAxis = Math.abs(dx) > Math.abs(dy) ? "x" : "y";
    if (touchAxis === "y") {
      dragging.value = false;
      resume();
      return;
    }
  }
  event.preventDefault();
  moveDrag(touch.clientX);
}
function onTouchEnd() {
  endDrag();
}

// Lightbox — index into `props.photos` (never into the cloned `slides`).
const lightboxIndex = ref<number | null>(null);
const lightboxOpen = computed(() => lightboxIndex.value !== null);
const lightboxPanelRef = ref<HTMLElement | null>(null);
const { activate: activateTrap, deactivate: deactivateTrap } = useFocusTrap(lightboxPanelRef, {
  immediate: false,
  escapeDeactivates: true,
  onDeactivate: () => closeLightbox(),
});

// Fast next()/prev() clicking used to hit a blank frame while the freshly-
// requested photo was still downloading — the <img src> switched instantly
// but had nothing to paint yet. `lightboxSrc` only updates once a photo has
// actually finished loading, so the previously-shown photo stays on screen
// (dimmed, with a spinner over it) instead of going blank; `loadToken`
// discards any load whose result arrives after a newer one was already
// requested, so mashing next/prev can't make an older photo "win" and
// flash in out of order. Neighbors are warmed into the browser cache as
// soon as a photo is shown, so by the time someone actually clicks
// next/prev again the image is very likely already local.
const lightboxSrc = ref("");
const lightboxLoading = ref(false);
let loadToken = 0;

function preload(url: string) {
  const img = new Image();
  img.src = url;
}

function loadLightboxPhoto(i: number) {
  const url = resolveMediaUrl(props.photos[i]!.image);
  const token = ++loadToken;
  lightboxLoading.value = true;
  const img = new Image();
  img.onload = img.onerror = () => {
    if (token !== loadToken) return;
    if (img.complete && img.naturalWidth > 0) lightboxSrc.value = url;
    lightboxLoading.value = false;
  };
  img.src = url;

  if (props.photos.length > 1) {
    preload(resolveMediaUrl(props.photos[(i + 1) % props.photos.length]!.image));
    preload(
      resolveMediaUrl(props.photos[(i - 1 + props.photos.length) % props.photos.length]!.image),
    );
  }
}

// Toggled by clicking the lightbox photo itself — see the template. Reset
// on close/navigate so every photo opens at "fit to screen" first.
const zoomed = ref(false);

function openLightbox(photoIndex: number) {
  if (didDrag) {
    didDrag = false;
    return;
  }
  pause();
  lightboxIndex.value = photoIndex;
  loadLightboxPhoto(photoIndex);
}
function closeLightbox() {
  lightboxIndex.value = null;
  lightboxSrc.value = "";
  zoomed.value = false;
  resume();
}
// Guarded on `lightboxLoading` (not just left to the buttons' `:disabled`)
// so the ArrowLeft/ArrowRight keyboard path — which has no disabled state
// of its own — can't queue clicks faster than photos actually load either.
// Preloading already makes most of these near-instant, so in practice this
// is rarely felt; it only holds a click back on a slow connection or a
// genuinely fast flurry of clicks, exactly the cases that used to outrun
// the network and flash a stale-looking frame.
//
// `loadLightboxPhoto` is called right here rather than from a
// `watch(lightboxIndex, ...)` — a watcher only runs on the next reactive
// flush (a microtask later), so `lightboxLoading` wouldn't actually read as
// `true` yet to a second call landing in the same synchronous tick (e.g.
// two clicks dispatched back to back without letting the event loop turn in
// between). Setting it here makes the guard take effect immediately,
// exactly when it's checked.
// Same fixed-cooldown idea as the main carousel's arrows (see
// NAV_COOLDOWN_MS above) — on top of the `lightboxLoading` guard, not
// instead of it: loading only holds a click back when the network is
// actually slow, whereas this keeps clicks from firing faster than the
// fade transition even when every photo is already cached and loads
// instantly.
const lightboxCooldown = ref(false);
function startLightboxCooldown() {
  lightboxCooldown.value = true;
  setTimeout(() => {
    lightboxCooldown.value = false;
  }, NAV_COOLDOWN_MS);
}

function lightboxNext() {
  if (lightboxIndex.value === null || lightboxLoading.value || lightboxCooldown.value) return;
  const next = (lightboxIndex.value + 1) % props.photos.length;
  lightboxIndex.value = next;
  zoomed.value = false;
  loadLightboxPhoto(next);
  startLightboxCooldown();
}
function lightboxPrev() {
  if (lightboxIndex.value === null || lightboxLoading.value || lightboxCooldown.value) return;
  const prev = (lightboxIndex.value - 1 + props.photos.length) % props.photos.length;
  lightboxIndex.value = prev;
  zoomed.value = false;
  loadLightboxPhoto(prev);
  startLightboxCooldown();
}
function onLightboxKeydown(event: KeyboardEvent) {
  if (event.key === "ArrowLeft") lightboxPrev();
  else if (event.key === "ArrowRight") lightboxNext();
}

watch(lightboxOpen, async (open) => {
  if (open) {
    document.body.style.overflow = "hidden";
    await nextTick();
    activateTrap();
  } else {
    deactivateTrap();
    document.body.style.overflow = "";
  }
});

onMounted(() => {
  measureStep();
  resizeObserver = new ResizeObserver(measureStep);
  if (trackRef.value) resizeObserver.observe(trackRef.value);
  if (viewportRef.value) resizeObserver.observe(viewportRef.value);
  startAutoplay();
});
onUnmounted(() => {
  stopAutoplay();
  resizeObserver?.disconnect();
  if (snapTimer !== null) clearTimeout(snapTimer);
  if (resetTimer !== null) clearTimeout(resetTimer);
  document.body.style.overflow = "";
  // In case the component unmounts mid-drag (e.g. navigating away while
  // the mouse button is still held).
  window.removeEventListener("mousemove", onMouseMove);
  window.removeEventListener("mouseup", onMouseUp);
});
</script>

<template>
  <div
    v-if="photos.length"
    class="relative w-full"
    role="region"
    aria-roledescription="carousel"
    aria-label="Фотогалерея школы"
    @mouseenter="pause"
    @mouseleave="resume"
    @focusin="pause"
    @focusout="resume"
    @keydown="onKeydown"
  >
    <div ref="viewportRef" class="overflow-hidden">
      <div
        ref="trackRef"
        data-no-orbit
        class="flex touch-pan-y select-none gap-24 motion-reduce:transition-none"
        :class="[
          jumping || dragging ? '' : 'transition-transform duration-500 ease-out',
          dragging ? 'cursor-grabbing' : 'cursor-grab',
        ]"
        :style="{ transform: `translateX(-${displayOffsetPx}px)` }"
        @mousedown="onMouseDown"
        @touchstart="onTouchStart"
        @touchmove="onTouchMove"
        @touchend="onTouchEnd"
      >
        <button
          v-for="(photo, i) in slides"
          :key="`${photo.id}-${i}`"
          type="button"
          class="group aspect-[3/4] h-[40vh] w-auto shrink-0 overflow-hidden rounded-lg"
          :aria-label="`Открыть фото ${(i % photos.length) + 1} из ${photos.length} на весь экран`"
          @click="openLightbox(i % photos.length)"
        >
          <NuxtImg
            :src="thumbUrl(photo)"
            format="webp"
            :sizes="`${THUMBNAIL_WIDTH}px`"
            alt=""
            draggable="false"
            class="h-full w-full object-cover transition-transform duration-300 ease-out group-hover:scale-110 motion-reduce:transition-none motion-reduce:group-hover:scale-100"
            loading="lazy"
          />
        </button>
      </div>
    </div>

    <template v-if="loops">
      <button
        type="button"
        aria-label="Предыдущее фото"
        class="absolute left-8 top-1/2 grid size-[44px] -translate-y-1/2 place-items-center rounded-full bg-white/90 text-ink shadow-md transition-transform hover:enabled:scale-105 active:enabled:scale-95 disabled:opacity-40"
        :disabled="navCooldown"
        @click="onArrowClick(prev)"
      >
        <ChevronLeft class="size-24" aria-hidden="true" />
      </button>
      <button
        type="button"
        aria-label="Следующее фото"
        class="absolute right-8 top-1/2 grid size-[44px] -translate-y-1/2 place-items-center rounded-full bg-white/90 text-ink shadow-md transition-transform hover:enabled:scale-105 active:enabled:scale-95 disabled:opacity-40"
        :disabled="navCooldown"
        @click="onArrowClick(next)"
      >
        <ChevronRight class="size-24" aria-hidden="true" />
      </button>

      <div class="mt-16 px-16">
        <div
          class="flex w-full flex-wrap items-center justify-center gap-16 rounded-full bg-white/40 p-16 backdrop-blur-lg backdrop-saturate-150"
        >
          <button
            v-for="(photo, i) in photos"
            :key="photo.id"
            type="button"
            :aria-label="`Перейти к фото ${i + 1} из ${photos.length}`"
            :aria-current="i === activeDot"
            class="grid place-items-center p-4"
            @click="goTo(i)"
          >
            <span
              class="block size-8 rounded-full transition-colors"
              :class="i === activeDot ? 'bg-primary' : 'bg-primary/30'"
            />
          </button>
        </div>
      </div>
    </template>

    <Teleport to="body">
      <div
        v-if="lightboxOpen"
        ref="lightboxPanelRef"
        data-no-orbit
        class="fixed inset-0 z-50 flex items-center justify-center bg-ink/90 p-24"
        role="dialog"
        aria-modal="true"
        aria-label="Просмотр фото"
        @click.self="closeLightbox"
        @keydown="onLightboxKeydown"
      >
        <button
          type="button"
          aria-label="Закрыть"
          class="absolute right-16 top-16 grid size-[44px] place-items-center rounded-full bg-white/10 text-white hover:bg-white/20"
          @click="closeLightbox"
        >
          <X class="size-24" aria-hidden="true" />
        </button>

        <button
          v-if="photos.length > 1"
          type="button"
          aria-label="Предыдущее фото"
          class="absolute left-8 top-1/2 grid size-[44px] -translate-y-1/2 place-items-center rounded-full bg-white/10 text-white hover:enabled:bg-white/20 disabled:opacity-40 sm:left-24"
          :disabled="lightboxLoading || lightboxCooldown"
          @click="lightboxPrev"
        >
          <ChevronLeft class="size-24" aria-hidden="true" />
        </button>
        <button
          v-if="photos.length > 1"
          type="button"
          aria-label="Следующее фото"
          class="absolute right-8 top-1/2 grid size-[44px] -translate-y-1/2 place-items-center rounded-full bg-white/10 text-white hover:enabled:bg-white/20 disabled:opacity-40 sm:right-24"
          :disabled="lightboxLoading || lightboxCooldown"
          @click="lightboxNext"
        >
          <ChevronRight class="size-24" aria-hidden="true" />
        </button>

        <!--
          This wrapper (and the image inside it) must stay unpositioned
          (no `relative`/`absolute`) — a positioned box sized close to the
          image extends close enough to the viewport edges to shift into
          its own paint layer and intercept clicks meant for the close/
          prev/next buttons above (confirmed live: a real click on the
          close button landed on a wrapper like this once already). Left
          unpositioned, it paints below the (positioned) buttons regardless
          of geometric overlap, so it can never steal their clicks — the
          spinner below anchors straight to the dialog instead, which is
          already `position: fixed`.

          Not zoomed: centered, capped to fit the screen (90vh tall). Click
          toggles `zoomed`, which drops the cap so the image renders at its
          real (upload-capped-at-2000px) resolution — on most screens
          that's bigger than the viewport, so the wrapper switches from
          centering to `overflow-auto` and lets the browser's native
          scrollbars pan around it, same as zooming into any oversized
          image.
        -->
        <div
          class="max-h-[94vh] max-w-[94vw]"
          :class="zoomed ? 'overflow-auto' : 'flex items-center justify-center'"
        >
          <img
            v-if="lightboxSrc"
            :src="lightboxSrc"
            alt=""
            class="rounded-lg transition-opacity duration-150"
            :class="[
              lightboxLoading ? 'opacity-40' : 'opacity-100',
              zoomed
                ? 'max-w-none cursor-zoom-out'
                : 'h-[90vh] w-auto max-w-full object-contain cursor-zoom-in',
            ]"
            @click="zoomed = !zoomed"
          />
        </div>
        <svg
          v-if="lightboxLoading"
          class="pointer-events-none absolute left-1/2 top-1/2 size-48 -translate-x-1/2 -translate-y-1/2 animate-spin text-white motion-reduce:hidden"
          viewBox="0 0 24 24"
          fill="none"
          aria-hidden="true"
        >
          <circle
            class="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            stroke-width="4"
          />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8v4a4 4 0 0 0-4 4H4z" />
        </svg>
      </div>
    </Teleport>
  </div>
</template>
