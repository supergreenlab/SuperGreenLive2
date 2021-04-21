export default async ({ store }) => {
  addEventListener('load', async () => {
    try {
      const urlParams = new URLSearchParams(window.location.search)
      if (urlParams.get('promo')) {
        await store.dispatch('checkout/fetchPromocode', {code: urlParams.get('promo')})
      }
    } catch(e) {
      console.log(e)
    }
  })
}
