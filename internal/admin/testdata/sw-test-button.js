Component.register('sw-test-button', {
    props: {
        variant: {
            type: String,
            required: false,
            default: 'primary'
        },
        size: {
            type: String,
            required: false,
            default: 'medium'
        },
        disabled: {
            type: Boolean,
            required: false,
            default: false
        }
    },

    computed: {
        buttonClasses() {
            return {
                [`sw-button--${this.variant}`]: this.variant,
                [`sw-button--${this.size}`]: this.size
            };
        },
        
        isDisabled() {
            return this.disabled;
        }
    },

    methods: {
        onClick() {
            this.$emit('click');
        }
    }
});
