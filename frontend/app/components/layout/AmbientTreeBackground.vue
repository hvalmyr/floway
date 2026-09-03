<script setup lang="ts">
import * as THREE from "three";
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls.js";
import { onBeforeUnmount, onMounted, ref } from "vue";
import {
  buildFloweryBranch,
  disposeObject3D,
  isReasonableBranch,
} from "~/lib/three/flowery-branch";

// Fixed, transparent 3D backdrop behind every page — the branch sits in a
// void behind the "glass" of the UI. The canvas itself never receives real
// pointer events (pointer-events: none, inherited from the wrapper below —
// nothing on the page is ever blocked by it); window-level listeners
// re-dispatch gestures onto the canvas instead.
//
// Interaction is entirely gated behind an explicit "focus" state, on both
// desktop and mobile — nothing on the object responds until the user
// double-clicks/double-taps it (see enterFocus below), and everything
// stops responding again the moment they leave focus. Desktop, while
// focused: left-drag rotates (from anywhere on screen), right-drag pans,
// scroll zooms. Mobile, while focused: one finger rotates, two fingers
// pan/pinch-zoom. Leaving focus happens three ways: a double click/tap
// anywhere at all (see trackClickEnd — unlike entering, leaving isn't tied
// to hitting the object), a press on real site UI — a header/footer link,
// a button, a form field (see isFunctionalElement, used in
// handlePointerDown/handleTouchStart) — or a few seconds of no interaction.

const canvasRef = ref<HTMLCanvasElement | null>(null);
const flashRef = ref<HTMLDivElement | null>(null);

let renderer: THREE.WebGLRenderer | null = null;
let controls: OrbitControls | null = null;
let scene: THREE.Scene | null = null;
let camera: THREE.PerspectiveCamera | null = null;
let branch: THREE.Group | null = null;
let branchSphere: THREE.Sphere | null = null;
let activationRing: THREE.Mesh | null = null;
let resizeObserver: ResizeObserver | null = null;
let frameId = 0;
let disposed = false;
let autoRotateResumeTimer: number | null = null;
let reducedMotion = false;

// Pauses auto-rotate the instant the user grabs the object (or enters
// touch focus mode), resumes 10s after they let go — shared by the mouse
// drag-release path (controls' own "end" event) and the touch focus-exit
// path (exitFocus below) so both go through identical scheduling.
function scheduleAutoRotateResume() {
  if (reducedMotion) return;
  if (autoRotateResumeTimer !== null) window.clearTimeout(autoRotateResumeTimer);
  autoRotateResumeTimer = window.setTimeout(() => {
    autoRotateResumeTimer = null;
    if (controls) controls.autoRotate = true;
  }, 10_000);
}

// One unified focus state for both platforms: nothing on the object
// responds to input until the user double-clicks/double-taps it (see the
// click-tracking further down), and a few seconds of no interaction exits
// it automatically, same as an explicit exit press does.
const FOCUS_IDLE_MS = 3000;
let focusActive = false;
let focusIdleTimer: number | null = null;

// Against the actual mesh, a sparse asymmetric branch (thin twigs, small
// flowers) has a tiny real hit area relative to how big it reads visually —
// missing it on every other attempt, worse still with an imprecise touch.
// Testing against its bounding sphere instead makes double-click/double-tap
// forgiving: anywhere that reads as "on the tree" to the eye now counts.
// Only relevant while *entering* focus — once focused, a double click/tap
// leaves focus regardless of where it lands (see trackClickEnd), so there's
// no equivalent "hit test" needed for exiting.
function hitsBranchLoose(clientX: number, clientY: number): boolean {
  if (!camera || !branchSphere) return false;
  const ndc = new THREE.Vector2(
    (clientX / window.innerWidth) * 2 - 1,
    -(clientY / window.innerHeight) * 2 + 1,
  );
  const raycaster = new THREE.Raycaster();
  raycaster.setFromCamera(ndc, camera);
  return raycaster.ray.intersectsSphere(branchSphere);
}

// preventDefault() on our pointermove handler alone doesn't reliably stop
// native scroll/pinch-zoom: the browser's compositor decides whether a
// touch pans/zooms the page from the `touch-action` in effect at that
// touch's own touchstart, before our JS ever runs — a later preventDefault
// is too late. Toggling touch-action:none on the document root (not just
// the canvas, which is pointer-events:none and never the real touch
// target) closes that gap for the whole time focus mode is active.
//
// Same root issue for text selection and image drag: a press held on
// running text or an image starts the browser's own native long-press
// selection/drag, and neither touch-action nor a plain preventDefault
// covers that — user-select and -webkit-user-drag are the dedicated CSS
// switches for it, toggled the same way.
let focusInteractionLocked = false;
let previousTouchAction = "";
let previousUserSelect = "";
let previousWebkitUserSelect = "";
let previousWebkitUserDrag = "";

function lockFocusInteraction() {
  if (focusInteractionLocked) return;
  focusInteractionLocked = true;
  const style = document.documentElement.style;
  previousTouchAction = style.touchAction;
  previousUserSelect = style.userSelect;
  previousWebkitUserSelect = style.getPropertyValue("-webkit-user-select");
  previousWebkitUserDrag = style.getPropertyValue("-webkit-user-drag");
  style.touchAction = "none";
  style.userSelect = "none";
  style.setProperty("-webkit-user-select", "none");
  style.setProperty("-webkit-user-drag", "none");
}

function unlockFocusInteraction() {
  if (!focusInteractionLocked) return;
  focusInteractionLocked = false;
  const style = document.documentElement.style;
  style.touchAction = previousTouchAction;
  style.userSelect = previousUserSelect;
  style.setProperty("-webkit-user-select", previousWebkitUserSelect);
  style.setProperty("-webkit-user-drag", previousWebkitUserDrag);
}

function enterFocus() {
  if (focusActive) return;
  focusActive = true;
  if (controls) controls.autoRotate = false;
  if (autoRotateResumeTimer !== null) {
    window.clearTimeout(autoRotateResumeTimer);
    autoRotateResumeTimer = null;
  }
  lockFocusInteraction();
  triggerFlash();
  playFocusSound();
  playFocusRingAnim();
  armFocusIdleTimer();
}

function exitFocus() {
  if (!focusActive) return;
  focusActive = false;
  if (focusIdleTimer !== null) {
    window.clearTimeout(focusIdleTimer);
    focusIdleTimer = null;
  }
  gestureAnchor = null;
  lastPinchDistance = null;
  unlockFocusInteraction();
  triggerFlash();
  playFocusSound();
  playDefocusRingAnim();
  scheduleAutoRotateResume();
}

function armFocusIdleTimer() {
  if (focusIdleTimer !== null) window.clearTimeout(focusIdleTimer);
  focusIdleTimer = window.setTimeout(exitFocus, FOCUS_IDLE_MS);
}

// A quick full-screen flash (not a sustained overlay) marks the moment
// focus changes — entering and leaving alike (double-click/double-tap
// either way, or the idle timeout) — restarting the CSS animation via a
// forced reflow so back-to-back triggers each get their own visible pulse.
function triggerFlash() {
  if (reducedMotion) return;
  const el = flashRef.value;
  if (!el) return;
  el.classList.remove("ambient-tree-flash-active");
  void el.offsetWidth;
  el.classList.add("ambient-tree-flash-active");
}

// Same moment, same cue, just audible — a single shared clip for both
// entering and leaving focus. Rewinding to the start on every call (rather
// than creating a new Audio each time) lets a rapid enter-then-exit
// retrigger it cleanly instead of the two overlapping.
const FOCUS_SOUND_SRC = "/sounds/focus-toggle.wav";
let focusSound: HTMLAudioElement | null = null;

function playFocusSound() {
  if (!focusSound) return;
  focusSound.currentTime = 0;
  // Playback can be rejected (e.g. a browser autoplay quirk) — this is
  // always triggered by a real click/tap, so that shouldn't happen, but a
  // rejected promise would otherwise surface as an unhandled rejection.
  void focusSound.play().catch(() => {});
}

// The ring is otherwise idle/invisible — these just start a one-shot tween
// that the render loop below advances every frame. Focusing: the ring
// appears and its radius shrinks to nothing, reading as "converging onto
// the object". Defocusing: the reverse motion (radius grows from nothing)
// while fading out, reading as "dissolving away" by the time it's done.
type RingAnimPhase = "idle" | "focusing" | "defocusing";
const RING_ANIM_MS = 500;
let ringAnimPhase: RingAnimPhase = "idle";
let ringAnimStart = 0;

function playFocusRingAnim() {
  if (reducedMotion || !activationRing) return;
  ringAnimPhase = "focusing";
  ringAnimStart = performance.now();
  activationRing.visible = true;
}

function playDefocusRingAnim() {
  if (reducedMotion || !activationRing) return;
  ringAnimPhase = "defocusing";
  ringAnimStart = performance.now();
  activationRing.visible = true;
}

function distance(a: { x: number; y: number }, b: { x: number; y: number }): number {
  return Math.hypot(a.x - b.x, a.y - b.y);
}

// Double-click is a page-wide gesture now (see enterFocus/exitFocus below),
// so its native side effect — select-word, select-image — would otherwise
// fire right alongside it almost anywhere, even where the cursor only
// *looks* like it's over empty space: a paragraph's line box, or an
// image's own box, extends past its visible glyphs/pixels, so a
// double-click landing in that gap still selects the nearest word/image.
// `detail` is the browser's own click-count (2 for the second click of a
// double-click) — checking it here, before the browser acts on it, is what
// actually suppresses the selection; a later 'dblclick' listener would
// already be too late, since selection happens as part of this same event's
// default action. Listens on 'mousedown', not 'pointerdown': Chromium never
// populates PointerEvent.detail (it's always 0), so the click count is only
// readable off the compatibility MouseEvent that follows it — 'mousedown'
// still fires in time to preventDefault() before selection happens. Left
// alone inside genuinely editable fields, where double-click-to-select-word
// is still useful.
const EDITABLE_SELECTOR = 'input, textarea, select, [contenteditable="true"]';

function suppressNativeDoubleClick(e: MouseEvent) {
  if (e.detail < 2) return;
  if (e.target instanceof Element && e.target.closest(EDITABLE_SELECTOR)) return;
  e.preventDefault();
}

// Entirely separate from the focused-interaction forwarding below: watches
// every mouse/touch press-release cycle purely to classify clicks/taps, so
// it can tell a double-click/double-tap apart from an ordinary one — on the
// branch to *enter* focus (see isOrbitExcluded/hitsBranchLoose below), or
// anywhere at all, while already focused, to *leave*. The same
// short-duration, small-movement test that keeps a real rotate/pan drag
// from ever being misread as a "double-click to exit" applies equally to
// entering, so leaving needed nothing new beyond remembering which of the
// two directions a pending click belongs to (see the `exit` flag below) —
// a stray entry-click and a stray exit-click shouldn't pair up with each
// other.
const CLICK_MAX_DURATION_MS = 300;
const CLICK_MAX_MOVE_PX = 12;
const DOUBLE_CLICK_WINDOW_MS = 350;
const DOUBLE_CLICK_MAX_DISTANCE_PX = 40;

const pressStartInfo = new Map<number, { x: number; y: number; time: number }>();
let pendingClick: { x: number; y: number; time: number; exit: boolean } | null = null;
let pendingClickTimer: number | null = null;

function clearPendingClick() {
  pendingClick = null;
  if (pendingClickTimer !== null) {
    window.clearTimeout(pendingClickTimer);
    pendingClickTimer = null;
  }
}

// Touch, or the mouse's primary (left) button only — a double right-click
// shouldn't enter/leave focus, since right-drag is reserved for panning
// once already focused.
function isEntryPointerType(e: PointerEvent): boolean {
  return e.pointerType === "touch" || (e.pointerType === "mouse" && e.button === 0);
}

// isTrusted guards throughout: the pointerup forwarded by endTouch bubbles
// (see its comment for why) and would otherwise reach this listener too,
// getting misread as a genuine click/tap release.
function trackClickStart(e: PointerEvent) {
  if (!e.isTrusted || !isEntryPointerType(e)) return;
  pressStartInfo.set(e.pointerId, { x: e.clientX, y: e.clientY, time: performance.now() });
}

function trackClickCancel(e: PointerEvent) {
  if (!e.isTrusted || !isEntryPointerType(e)) return;
  pressStartInfo.delete(e.pointerId);
}

function trackClickEnd(e: PointerEvent) {
  if (!e.isTrusted || !isEntryPointerType(e)) return;
  const start = pressStartInfo.get(e.pointerId);
  pressStartInfo.delete(e.pointerId);
  if (!start) return;

  const end = { x: e.clientX, y: e.clientY };
  const duration = performance.now() - start.time;
  if (duration > CLICK_MAX_DURATION_MS || distance(start, end) > CLICK_MAX_MOVE_PX) {
    // A drag, not a click — breaks any pending double-click sequence.
    clearPendingClick();
    return;
  }

  // While already focused, a double click/tap leaves focus no matter where
  // it lands — entering is the only direction that needs to land on the
  // object (loosely) and skip real site controls, so a double-click meant
  // for a button/link/word-select never also enters focus by accident.
  if (!focusActive && (isOrbitExcluded(e.target) || !hitsBranchLoose(end.x, end.y))) {
    clearPendingClick();
    return;
  }

  const click = { x: end.x, y: end.y, time: performance.now(), exit: focusActive };
  if (
    pendingClick &&
    pendingClick.exit === click.exit &&
    click.time - pendingClick.time <= DOUBLE_CLICK_WINDOW_MS &&
    distance(pendingClick, click) <= DOUBLE_CLICK_MAX_DISTANCE_PX
  ) {
    clearPendingClick();
    if (click.exit) exitFocus();
    else enterFocus();
  } else {
    clearPendingClick();
    pendingClick = click;
    pendingClickTimer = window.setTimeout(clearPendingClick, DOUBLE_CLICK_WINDOW_MS);
  }
}

// Real, interactive site controls — as opposed to plain text, which reads
// fine but isn't something you *do* anything with. Split in two: pressing
// either kind while entering focus doesn't count as the double-click that
// enters it (isOrbitExcluded, below); pressing a functional one specifically
// while *already* focused exits focus outright and lets the press through
// untouched (see handlePointerDown/handleTouchStart further down) — using
// the site's real UI (a header link, a form field, anything with its own
// job to do) should never quietly get eaten as a rotate/pan gesture instead.
const FUNCTIONAL_SELECTOR =
  'a, button, input, textarea, select, [role="button"], [contenteditable="true"], label, [data-no-orbit]';
const TEXT_SELECTOR =
  "p, h1, h2, h3, h4, h5, h6, span, li, blockquote, td, th, dt, dd, figcaption, strong, em";
const ORBIT_EXCLUDE_SELECTOR = `${FUNCTIONAL_SELECTOR}, ${TEXT_SELECTOR}`;

function isOrbitExcluded(target: EventTarget | null): boolean {
  return target instanceof Element && target.closest(ORBIT_EXCLUDE_SELECTOR) !== null;
}

function isFunctionalElement(target: EventTarget | null): boolean {
  return target instanceof Element && target.closest(FUNCTIONAL_SELECTOR) !== null;
}

// Forwarded events are built explicitly (not by passing the real event as
// the init dict) and forced non-bubbling: canvas is a descendant of
// `window`, so a *bubbling* synthetic event dispatched on it would climb
// back up to these same window listeners and re-trigger itself forever.

// Mouse only ever drives the object while focused — nothing responds to a
// plain click/drag otherwise. Left drags rotate from anywhere on screen
// once focused, not just from the object itself (standard orbit-control
// feel) — leaving focus is its own explicit double-click gesture (see
// trackClickEnd), not tied to where a drag happens to start or end.
//
// A press on real site UI is the other way to leave focus, and it's
// unconditional (no double-press needed): using the actual site — a
// header/footer link, a button, a form field — should never be swallowed
// as a rotate/pan attempt just because focus happened to still be active.
// Exits and returns immediately, without forwarding or preventDefault, so
// the press behaves exactly like it would if focus had never existed.
function handlePointerDown(e: PointerEvent) {
  const canvas = canvasRef.value;
  if (!canvas || e.pointerType !== "mouse" || !focusActive) return;
  if (isFunctionalElement(e.target)) {
    exitFocus();
    return;
  }
  e.preventDefault();
  canvas.dispatchEvent(
    new PointerEvent("pointerdown", {
      bubbles: false,
      cancelable: true,
      pointerId: e.pointerId,
      pointerType: e.pointerType,
      clientX: e.clientX,
      clientY: e.clientY,
      button: e.button,
      buttons: e.buttons,
      ctrlKey: e.ctrlKey,
      shiftKey: e.shiftKey,
      altKey: e.altKey,
      metaKey: e.metaKey,
      isPrimary: e.isPrimary,
    }),
  );
}

// Scroll only zooms the object while focused — reading the page never
// nudges it by accident. Needs a non-passive listener since the real event
// is prevented once active, so zooming doesn't also scroll the page
// underneath it.
function handleWheel(e: WheelEvent) {
  const canvas = canvasRef.value;
  if (!canvas || !focusActive) return;
  e.preventDefault();
  canvas.dispatchEvent(
    new WheelEvent("wheel", {
      bubbles: false,
      cancelable: true,
      clientX: e.clientX,
      clientY: e.clientY,
      deltaX: e.deltaX,
      deltaY: e.deltaY,
      deltaZ: e.deltaZ,
      deltaMode: e.deltaMode,
      ctrlKey: e.ctrlKey,
    }),
  );
}

// Right-drag-to-pan needs the browser's own context menu out of the way —
// only relevant while focused, since that's the only time right-drag does
// anything.
function handleContextMenu(e: MouseEvent) {
  if (!focusActive) return;
  e.preventDefault();
}

// Touch never drives the object unless explicitly focused (double-tap the
// branch first) — an unfocused finger, one or two, is left completely
// alone so the page scrolls/pinch-zooms exactly like any other page.
// Once focused, EVERY touch is forwarded regardless of what element it
// started on — even running text, a link, or the photo carousel — because
// that's the whole point of focus mode: the entire gesture space belongs
// to the object until the user explicitly leaves focus. Excluding some
// elements even while focused was the bug (a drag that happened to start
// or cross over text/a button would silently stop being tracked, since
// the exclusion check ran on every event, not just the one that starts
// the gesture).
const activeTouches = new Map<number, { x: number; y: number }>();

function touchPointerInit(
  pointerId: number,
  x: number,
  y: number,
  bubbles = false,
): PointerEventInit {
  return { bubbles, cancelable: true, pointerId, pointerType: "touch", clientX: x, clientY: y };
}

// Seeds OrbitControls' internal pointer tracking. Once seeded, real
// subsequent touchmove/touchup events would normally bubble to the
// document naturally and OrbitControls' own document-level listeners
// (added internally on the first forwarded pointerdown) would pick them
// up by matching pointerId — but handleFocusedGestureMove below
// intercepts all of them first (see its comment for why).
//
// One finger rotates from anywhere on screen once focused, not just from
// the object itself — leaving focus is its own explicit double-tap gesture
// (see trackClickEnd), not tied to where a drag happens to start or end.
//
// A tap on real site UI is the other way to leave focus — see
// handlePointerDown's comment for why. Only checked for the first finger
// of a fresh gesture (activeTouches still empty); a second finger joining
// an ongoing pinch/pan is never a fresh press in its own right.
function handleTouchStart(e: PointerEvent) {
  const canvas = canvasRef.value;
  if (!canvas || e.pointerType !== "touch" || !focusActive) return;
  if (activeTouches.size === 0 && isFunctionalElement(e.target)) {
    exitFocus();
    return;
  }
  activeTouches.set(e.pointerId, { x: e.clientX, y: e.clientY });
  gestureAnchor = null;
  lastPinchDistance = null;
  e.preventDefault();
  canvas.dispatchEvent(
    new PointerEvent("pointerdown", touchPointerInit(e.pointerId, e.clientX, e.clientY)),
  );
}

// A real held finger (or two) is never perfectly stationary — a pixel or
// two of sensor/skin-settle noise is normal — but OrbitControls has no
// built-in tolerance for that: every reported pixel of movement becomes
// rotation or pan, which read as the camera slowly drifting on its own
// even with the fingers deliberately held still. Dead-zoned against the
// *average* position of every active touch (rotate for one finger, pan
// for two, matching how OrbitControls itself derives _rotateStart/
// _panStart) — with an escape hatch so a real pinch still zooms
// immediately even while that average barely moves: if the distance
// *between* two fingers changes past the same threshold, that's forwarded
// regardless of the average-position dead-zone.
//
// Registered on the CAPTURE phase (unlike the rest of this file's
// listeners) specifically so it runs, and can stopPropagation(), before
// OrbitControls' own document-level listener — added on `document`, which
// is earlier than `window` in the bubble phase — ever sees the raw,
// unfiltered movement.
const GESTURE_DEAD_ZONE_PX = 4;
let gestureAnchor: { x: number; y: number } | null = null;
let lastPinchDistance: number | null = null;

function handleFocusedGestureMove(e: PointerEvent) {
  // Ignore our own forwarded moves (see below) — synthetic events are
  // never `isTrusted`, so this only ever runs for genuine touches.
  if (!e.isTrusted || e.pointerType !== "touch" || !focusActive || !activeTouches.has(e.pointerId))
    return;
  e.preventDefault();
  e.stopPropagation();
  activeTouches.set(e.pointerId, { x: e.clientX, y: e.clientY });

  const touches = [...activeTouches.entries()];
  const avg = {
    x: touches.reduce((sum, [, p]) => sum + p.x, 0) / touches.length,
    y: touches.reduce((sum, [, p]) => sum + p.y, 0) / touches.length,
  };
  const pinchDistance = touches.length === 2 ? distance(touches[0][1], touches[1][1]) : null;

  if (!gestureAnchor) {
    gestureAnchor = avg;
    lastPinchDistance = pinchDistance;
    return;
  }

  const avgMoved = distance(gestureAnchor, avg) >= GESTURE_DEAD_ZONE_PX;
  const pinchChanged =
    pinchDistance !== null &&
    lastPinchDistance !== null &&
    Math.abs(pinchDistance - lastPinchDistance) >= GESTURE_DEAD_ZONE_PX;
  if (!avgMoved && !pinchChanged) return;

  gestureAnchor = avg;
  lastPinchDistance = pinchDistance;
  const canvas = canvasRef.value;
  if (!canvas) return;
  // Every active touch is re-sent, not just the one that triggered this
  // event: OrbitControls tracks each pointerId's last known position
  // independently (for pan/pinch math involving both fingers), so the
  // finger that *didn't* just cross the dead zone still needs its current
  // position on record.
  //
  // Dispatched with bubbles:true, unlike the other forwarded events in
  // this file: once the first finger's pointerdown is captured (via
  // setPointerCapture, inside OrbitControls' own handler), the *browser*
  // retargets every further event carrying that pointerId to the
  // capturing element (canvas) before dispatch — regardless of which
  // element .dispatchEvent() was actually called on. OrbitControls listens
  // for pointermove on canvas's *document*, an ancestor of canvas, so a
  // non-bubbling event silently never arrives there once retargeted. It
  // does climb back up through this same window-level listener, but the
  // isTrusted guard above makes that a no-op.
  for (const [id, pos] of touches) {
    canvas.dispatchEvent(new PointerEvent("pointermove", touchPointerInit(id, pos.x, pos.y, true)));
  }
}

// Forwards the lift regardless of the current focusActive value — if focus
// exited mid-gesture (idle timeout, tap elsewhere) while a finger was still
// down, OrbitControls is still tracking that pointerId from the forwarded
// pointerdown and needs the matching pointerup to release it cleanly.
//
// Dispatched with bubbles:true for the same reason as the move forward
// above — OrbitControls' pointerup listener is on canvas's document, not
// canvas itself, and once that pointerId is captured, its effective
// target is forced to canvas regardless of dispatch site. Safe to bubble
// back through this same listener: activeTouches.delete() already
// happened, so a re-entrant call for the same pointerId is a no-op.
function endTouch(e: PointerEvent) {
  const canvas = canvasRef.value;
  if (!canvas || e.pointerType !== "touch" || !activeTouches.delete(e.pointerId)) return;
  // Re-anchor from scratch on the next move, whether that's a full lift or
  // just dropping from two fingers to one — otherwise the remaining
  // finger's next move would be dead-zoned against a stale two-finger
  // average instead of its own actual position.
  gestureAnchor = null;
  lastPinchDistance = null;
  canvas.dispatchEvent(
    new PointerEvent("pointerup", touchPointerInit(e.pointerId, e.clientX, e.clientY, true)),
  );
}

function disposeScene() {
  if (scene) disposeObject3D(scene);
}

// A fresh composition on every load, but not an ugly one: reject a bad
// roll and try another seed rather than accept whatever came out, with
// the original hand-picked seed as an always-safe last resort.
function pickBranch(): THREE.Group {
  for (let attempt = 0; attempt < 8; attempt++) {
    const seed = Math.floor(Math.random() * 0xffffffff);
    const candidate = buildFloweryBranch(seed);
    if (isReasonableBranch(candidate)) return candidate;
    disposeObject3D(candidate);
  }
  return buildFloweryBranch(2024);
}

// Bounding-box/sphere center often lands in visual empty space for a
// sparse, asymmetric shape like a branch — orbiting around it reads as
// circling an arbitrary point in the void rather than the object itself.
// A vertex-density-weighted centroid sits inside the actual mass instead.
function computeCentroid(object: THREE.Object3D): THREE.Vector3 {
  const centroid = new THREE.Vector3();
  const v = new THREE.Vector3();
  let count = 0;
  object.traverse((child) => {
    if (!(child instanceof THREE.Mesh)) return;
    const pos = child.geometry.attributes.position;
    if (!pos) return;
    for (let i = 0; i < pos.count; i++) {
      v.fromBufferAttribute(pos, i).applyMatrix4(child.matrixWorld);
      centroid.add(v);
      count++;
    }
  });
  if (count > 0) centroid.divideScalar(count);
  return centroid;
}

onMounted(() => {
  const canvas = canvasRef.value;
  const parent = canvas?.parentElement;
  if (!canvas || !parent) return;

  focusSound = new Audio(FOCUS_SOUND_SRC);
  focusSound.volume = 0.5;
  focusSound.preload = "auto";

  scene = new THREE.Scene();
  camera = new THREE.PerspectiveCamera(40, 1, 0.01, 100);

  branch = pickBranch();
  // Off-center yaw so it reads as something glimpsed to the side rather
  // than a centered hero object.
  branch.rotation.y = Math.PI * 0.15;
  branch.traverse((o) => {
    if ((o as THREE.Mesh).isMesh) {
      (o as THREE.Mesh).castShadow = true;
      (o as THREE.Mesh).receiveShadow = true;
    }
  });
  scene.add(branch);
  branch.updateMatrixWorld(true);

  const box = new THREE.Box3().setFromObject(branch);
  const sphere = box.getBoundingSphere(new THREE.Sphere());
  branchSphere = sphere.clone();
  const centroid = computeCentroid(branch);
  const dist = (sphere.radius / Math.tan((camera.fov * Math.PI) / 360)) * 1.7;
  const dir = new THREE.Vector3(0.9, 0.5, 1.1).normalize();
  camera.position.copy(centroid).addScaledVector(dir, dist);
  camera.near = Math.max(dist / 100, 0.01);
  camera.far = dist * 100;
  camera.updateProjectionMatrix();

  // "You just focused/unfocused this" cue — a soft accent-colored ring at
  // the object's base, hidden except during its one-shot converge/dissolve
  // tween (see playFocusRingAnim/playDefocusRingAnim and the render loop
  // below). Built at a fixed canonical radius; the tween animates
  // ring.scale rather than rebuilding geometry every frame.
  const ring = new THREE.Mesh(
    new THREE.RingGeometry(sphere.radius * 0.98, sphere.radius * 1.06, 48),
    new THREE.MeshBasicMaterial({
      color: 0x82b1cc,
      transparent: true,
      opacity: 1,
      side: THREE.DoubleSide,
      depthWrite: false,
    }),
  );
  ring.name = "activation_ring";
  ring.visible = false;
  ring.rotation.x = -Math.PI / 2;
  ring.position.set(centroid.x, box.min.y + 0.002, centroid.z);
  activationRing = ring;
  scene.add(ring);

  // Invisible except where the branch's shadow falls on it — keeps the
  // "floating in transparent space" look while still grounding the object.
  const groundSpan =
    Math.max(box.getSize(new THREE.Vector3()).x, box.getSize(new THREE.Vector3()).z) * 6;
  const ground = new THREE.Mesh(
    new THREE.PlaneGeometry(groundSpan, groundSpan),
    new THREE.ShadowMaterial({ opacity: 0.22 }),
  );
  ground.rotation.x = -Math.PI / 2;
  ground.position.y = box.min.y;
  ground.receiveShadow = true;
  scene.add(ground);

  scene.add(new THREE.HemisphereLight(0xffffff, 0xd8d2c4, 1.0));
  const key = new THREE.DirectionalLight(0xffffff, 1.8);
  key.position.set(4, 7, 5);
  key.castShadow = true;
  // 512 not 1024 — a soft, diffuse contact shadow doesn't need 1K
  // resolution, and the shadow pass cost scales with map area (this alone
  // cuts it to a quarter). See the render-cost comment by the renderer
  // below for the fuller picture.
  key.shadow.mapSize.set(512, 512);
  key.shadow.bias = -0.0003;
  const shadowSpan = sphere.radius * 2.2;
  key.shadow.camera.left = -shadowSpan;
  key.shadow.camera.right = shadowSpan;
  key.shadow.camera.top = shadowSpan;
  key.shadow.camera.bottom = -shadowSpan;
  key.shadow.camera.near = 0.01;
  key.shadow.camera.far = shadowSpan * 4;
  key.shadow.camera.updateProjectionMatrix();
  scene.add(key);
  const fill = new THREE.DirectionalLight(0xfff4e6, 0.4);
  fill.position.set(-5, 3, -4);
  scene.add(fill);

  // This renders continuously forever (see the render loop below —
  // autoRotate means there's always something to redraw), so per-frame
  // cost matters far more here than for a one-off render: a real mobile
  // Lighthouse audit traced the site's whole main-thread/CPU-idle failure
  // to this loop never letting the CPU rest. antialias off and a capped
  // pixel ratio are the two biggest per-frame levers — the scene sits at
  // 70% opacity behind blurred "glass" UI (see the template), so neither
  // is very visible here even though they would be on focal content.
  renderer = new THREE.WebGLRenderer({ canvas, alpha: true, antialias: false });
  renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 1.5));
  renderer.shadowMap.enabled = true;
  renderer.shadowMap.type = THREE.PCFSoftShadowMap;

  // enableZoom / enablePan default to true, and the default mouseButtons
  // mapping is already LEFT=ROTATE, MIDDLE=DOLLY, RIGHT=PAN — exactly the
  // "left rotates, right pans, scroll zooms" scheme requested. None of it
  // is reachable until focused, though — see handlePointerDown/handleWheel.
  controls = new OrbitControls(camera, canvas);
  controls.target.copy(centroid);
  controls.enableDamping = true;
  controls.dampingFactor = 0.08;
  controls.minDistance = dist * 0.35;
  controls.maxDistance = dist * 3;
  // Camera can't dip below the ground plane: this caps how far the tilt
  // can go toward the floor so, even at maxDistance, camera.y never drops
  // below box.min.y (with a 15% safety margin so it doesn't visually graze
  // the floor at the limit).
  const floorMarginY = Math.max(0.001, (centroid.y - box.min.y) * 0.85);
  const floorPolarLimit = Math.acos(
    THREE.MathUtils.clamp(-floorMarginY / controls.maxDistance, -1, 1),
  );
  // minPolarAngle defaulting to 0 (straight overhead) left ~70° of
  // "tilt toward looking from above" travel but only ~20° of "tilt toward
  // the floor" travel before hitting floorPolarLimit — that lopsided range
  // is exactly what made vertical drag feel far more sensitive in one
  // direction than the other (a small drag toward the floor immediately
  // hit the hard stop). Centering a fixed, symmetric cone of travel on the
  // camera's actual starting angle instead gives both directions the same
  // amount of room, still never exceeding the floor safety limit.
  const initialPolarAngle = new THREE.Spherical().setFromVector3(
    camera.position.clone().sub(centroid),
  ).phi;
  const verticalTiltRange = Math.PI / 4;
  controls.maxPolarAngle = Math.min(floorPolarLimit, initialPolarAngle + verticalTiltRange);
  // Clamped against the resolved maxPolarAngle (not just 0) so a branch
  // whose floorPolarLimit lands below initialPolarAngle - range still ends
  // up with a valid min <= max range instead of an inverted one.
  controls.minPolarAngle = Math.min(
    controls.maxPolarAngle,
    Math.max(0, initialPolarAngle - verticalTiltRange),
  );
  // Touch is only ever forwarded while focused (see handleTouchStart), so
  // this mapping only matters in focus mode: one finger rotates, two
  // fingers pinch-zoom/pan.
  controls.touches = { ONE: THREE.TOUCH.ROTATE, TWO: THREE.TOUCH.DOLLY_PAN };
  reducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  controls.autoRotate = !reducedMotion;
  controls.autoRotateSpeed = 0.5;
  // Pauses the instant the user grabs the object, resumes 10s after they
  // let go — cleared/rescheduled on every start/end so a resumed spin
  // doesn't jump back in mid-interaction if they grab it again quickly.
  if (!reducedMotion) {
    controls.addEventListener("start", () => {
      if (!controls) return;
      controls.autoRotate = false;
      if (autoRotateResumeTimer !== null) {
        window.clearTimeout(autoRotateResumeTimer);
        autoRotateResumeTimer = null;
      }
      // Any interaction while focused pushes its idle-exit window back out.
      if (focusActive) armFocusIdleTimer();
    });
    controls.addEventListener("end", scheduleAutoRotateResume);
  }
  controls.update();
  window.addEventListener("pointerdown", handlePointerDown);
  window.addEventListener("pointerdown", handleTouchStart);
  window.addEventListener("pointermove", handleFocusedGestureMove, { capture: true });
  window.addEventListener("pointerup", endTouch);
  window.addEventListener("pointercancel", endTouch);
  window.addEventListener("pointerdown", trackClickStart);
  window.addEventListener("pointerup", trackClickEnd);
  window.addEventListener("pointercancel", trackClickCancel);
  window.addEventListener("mousedown", suppressNativeDoubleClick);
  window.addEventListener("wheel", handleWheel, { passive: false });
  window.addEventListener("contextmenu", handleContextMenu);

  // Shifts where the (always object-centered-orbit) tree lands on screen,
  // fresh each load — setViewOffset re-frames the rendered window inside
  // the full frustum without touching the camera's actual position/target,
  // so orbiting still pivots on the object correctly regardless.
  const screenOffsetX = (Math.random() - 0.5) * 0.44;
  const screenOffsetY = (Math.random() - 0.5) * 0.24;

  const fit = () => {
    if (!renderer || !camera) return;
    const w = parent.clientWidth || window.innerWidth;
    const h = parent.clientHeight || window.innerHeight;
    if (!w || !h) return;
    renderer.setSize(w, h, false);
    camera.aspect = w / h;
    camera.setViewOffset(w, h, screenOffsetX * w, screenOffsetY * h, w, h);
    camera.updateProjectionMatrix();
  };
  fit();
  // Belt-and-suspenders: on the odd load where layout genuinely isn't
  // settled yet at mount time (both clientWidth and the innerWidth
  // fallback read 0 for that one synchronous instant), the guard above
  // skips sizing rather than leaving a 0×0 drawing buffer. Retried via
  // both rAF (next paint) and a plain timer — rAF alone doesn't fire until
  // the tab actually composites a frame, which a backgrounded/throttled
  // tab may delay far longer than a timer.
  requestAnimationFrame(fit);
  window.setTimeout(fit, 100);
  resizeObserver = new ResizeObserver(fit);
  resizeObserver.observe(parent);

  // The actual GPU draw (renderer.render, below) is what a real mobile
  // Lighthouse audit traced this component's whole CPU-idle failure to —
  // this runs forever (autoRotate means there's always something to
  // redraw), so at an uncapped rAF rate that's shadow-mapped draw calls
  // every ~8-16ms, indefinitely, on every page. A slow ambient rotation
  // doesn't need 60-120fps to read as smooth, so only the render call
  // itself is throttled to ~30fps — controls.update() (cheap vector math,
  // not a GPU cost) still runs every real frame so the damping/rotation
  // math stays exactly as smooth and doesn't visibly slow down.
  const RENDER_INTERVAL_MS = 1000 / 30;
  let lastRenderTime = 0;

  const loop = (now: number) => {
    if (disposed) return;
    frameId = requestAnimationFrame(loop);
    if (document.hidden || !controls || !renderer || !scene || !camera) return;
    controls.update();
    if (now - lastRenderTime < RENDER_INTERVAL_MS) return;
    lastRenderTime = now;
    if (activationRing && ringAnimPhase !== "idle") {
      const t = Math.min((performance.now() - ringAnimStart) / RING_ANIM_MS, 1);
      const mat = activationRing.material as THREE.MeshBasicMaterial;
      if (ringAnimPhase === "focusing") {
        // Ease-out: fast collapse at first, settling precisely onto zero —
        // paired with an opacity fade *in*, from fully transparent to
        // fully opaque, so it reads as solidifying into a point right as
        // it converges rather than just shrinking.
        const eased = 1 - (1 - t) ** 3;
        activationRing.scale.setScalar(Math.max(1 - eased, 0.001));
        mat.opacity = eased;
      } else {
        // Ease-in growth paired with a matching opacity fade, so it's
        // fully transparent right as it reaches full size — "dissolves".
        const eased = t * t;
        activationRing.scale.setScalar(Math.max(eased, 0.001));
        mat.opacity = 1 - t;
      }
      if (t >= 1) {
        ringAnimPhase = "idle";
        activationRing.visible = false;
      }
    }
    renderer.render(scene, camera);
  };
  loop(performance.now());
});

onBeforeUnmount(() => {
  disposed = true;
  cancelAnimationFrame(frameId);
  if (autoRotateResumeTimer !== null) window.clearTimeout(autoRotateResumeTimer);
  if (focusIdleTimer !== null) window.clearTimeout(focusIdleTimer);
  clearPendingClick();
  focusActive = false;
  gestureAnchor = null;
  lastPinchDistance = null;
  unlockFocusInteraction();
  resizeObserver?.disconnect();
  window.removeEventListener("pointerdown", handlePointerDown);
  window.removeEventListener("pointerdown", handleTouchStart);
  window.removeEventListener("pointermove", handleFocusedGestureMove, { capture: true });
  window.removeEventListener("pointerup", endTouch);
  window.removeEventListener("pointercancel", endTouch);
  window.removeEventListener("pointerdown", trackClickStart);
  window.removeEventListener("pointerup", trackClickEnd);
  window.removeEventListener("pointercancel", trackClickCancel);
  window.removeEventListener("mousedown", suppressNativeDoubleClick);
  window.removeEventListener("wheel", handleWheel);
  window.removeEventListener("contextmenu", handleContextMenu);
  activeTouches.clear();
  pressStartInfo.clear();
  focusSound?.pause();
  focusSound = null;
  controls?.dispose();
  disposeScene();
  renderer?.dispose();
});
</script>

<template>
  <div class="pointer-events-none fixed inset-0 -z-10 opacity-70" aria-hidden="true">
    <canvas ref="canvasRef" class="block h-full w-full" />
  </div>
  <!-- Whole-screen brief flash marking a focus change, entering or leaving
  — deliberately a transient pulse rather than a sustained dim/highlight
  overlay, and above everything (including the header) since it's a
  momentary UI-wide cue, not part of the 3D scene itself. -->
  <div
    ref="flashRef"
    class="pointer-events-none fixed inset-0 z-[100] bg-white opacity-0"
    aria-hidden="true"
  />
</template>

<style scoped>
.ambient-tree-flash-active {
  animation: ambient-tree-flash 450ms ease-out;
}

@keyframes ambient-tree-flash {
  0% {
    opacity: 0;
  }
  12% {
    opacity: 0.55;
  }
  100% {
    opacity: 0;
  }
}
</style>
