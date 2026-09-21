let injected = false;

/**
 * Injects the official Yandex.Metrika counter snippet (script + noscript
 * pixel), once per page load, client-side only. Called after the visitor
 * accepts the cookie banner (or immediately on mount for a returning visitor
 * who already accepted — see CookieConsentBanner.vue). No-ops without a
 * configured counter ID, e.g. local dev.
 */
export function useYandexMetrika() {
  const { public: publicConfig } = useRuntimeConfig();

  function load() {
    if (injected || import.meta.server) return;
    const counterId = Number(publicConfig.yandexMetrikaId);
    if (!counterId) return;
    injected = true;

    const script = document.createElement("script");
    script.textContent = `
      (function(m,e,t,r,i,k,a){m[i]=m[i]||function(){(m[i].a=m[i].a||[]).push(arguments)};
      m[i].l=1*new Date();
      for (var j = 0; j < document.scripts.length; j++) {if (document.scripts[j].src === r) { return; }}
      k=e.createElement(t),a=e.getElementsByTagName(t)[0],k.async=1,k.src=r,a.parentNode.insertBefore(k,a)})
      (window, document, "script", "https://mc.yandex.ru/metrika/tag.js", "ym");
      ym(${counterId}, "init", {
        clickmap: true,
        trackLinks: true,
        accurateTrackBounce: true,
        webvisor: true
      });
    `;
    document.head.appendChild(script);

    const noscript = document.createElement("noscript");
    const pixel = document.createElement("img");
    pixel.src = `https://mc.yandex.ru/watch/${counterId}`;
    pixel.style.cssText = "position:absolute;left:-9999px";
    pixel.alt = "";
    noscript.appendChild(pixel);
    document.body.appendChild(noscript);
  }

  return { load };
}
