const viewModules = import.meta.glob('../views/**/*.vue')
const placeholderView = viewModules['../views/placeholder/index.vue']()
const registeredRoutes = new Map()

function resolveView(node) {
  const key = `../views/${String(node).replace(/^\//, '')}/index.vue`
  const component = viewModules[key]
  return {
    component: component || (() => placeholderView),
    exists: Boolean(component),
  }
}

export function registerRouters(router, routers = []) {
  const visit = (items) => {
    items.forEach((item) => {
      if (item.types === 2 && item.route && !registeredRoutes.has(item.id)) {
        const view = resolveView(item.node)
        const remove = router.addRoute('Layout', {
          path: item.route.replace(/^\//, ''),
          name: item.node,
          component: view.component,
          meta: {
            routerId: item.id,
            title: item.title,
            icon: item.icon,
            keepAlive: view.exists && item.keep_alive === 1,
          },
        })
        registeredRoutes.set(item.id, remove)
      }
      if (item.children?.length) visit(item.children)
    })
  }
  visit(routers)
}

export function resetDynamicRoutes() {
  registeredRoutes.forEach((remove) => remove())
  registeredRoutes.clear()
}
