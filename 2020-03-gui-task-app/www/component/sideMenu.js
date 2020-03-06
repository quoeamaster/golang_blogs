
Vue.component('side-menu', {
    props: [ 'show' ],
    data: function () {
        return {
            'showContainer': false
        };
    },
    watch: {
        show: function (val) {
            this.showContainer = val;
        }
    },
    /*
    mounted: function() {
        alert(this.show);
        alert(this);
    },
    */
    methods: {
        /* ------------------ */
        /*   event handlers   */
        /* ------------------ */

        onCloseClick: function () {
            // console.log(this._events);
            this.$emit('update-show-side-menu', false);
        },


        /* --------------- */
        /*   css / style   */
        /* --------------- */

        getSideMenuClass: function () {
            // console.log('called get CLASS????' + this.show);
            if (this.showContainer === true) {
                return {
                    'side-menu-container-width': true
                };
            } else {
                return {
                    'side-menu-container-width': false
                };
            }
        },


    },
    template: `
<div class="side-menu-container"
    v-bind:class="getSideMenuClass()" >
    <!-- close button -->
    <div class="side-menu-close core-pointer" v-on:click="onCloseClick">&times;</div>
    <div class="side-menu-caption">menu</div>
    
    <div class="side-menu-item">option1</div>
    <div class="side-menu-item">about</div>
    
</div>
    `
});