Module.register('sw-product', {
    type: 'core',
    name: 'product',
    title: 'sw-product.general.mainMenuItemGeneral',
    description: 'sw-product.general.descriptionTextModule',
    color: '#57D9A3',
    icon: 'default-symbol-products',

    routes: {
        index: {
            component: 'sw-product-list',
            path: 'index'
        }
    },

    navigation: [{
        label: 'sw-product.general.mainMenuItemGeneral',
        color: '#57D9A3',
        path: 'sw.product.index',
        icon: 'default-symbol-products',
        parent: 'sw-catalogue',
        position: 10
    }]
});
