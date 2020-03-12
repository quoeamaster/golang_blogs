
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
    <div class="side-menu-caption">about</div>
    
    <div>
        <span class="side-menu-item" style="font-size: 1.4em;">
            daily notes app</span>
        <span class="side-menu-item" style="font-size: 1.0em;">
            version: 1.0.0</span>
        <span class="side-menu-item" style="font-size: 1.0em;">
            powered by: Que Master</span>
        <span class="side-menu-item" style="font-size: 1.0em;">
            <i class="far fa-envelope"></i> quoeamaster@gmail.com</span>
    </div>
</div>
    `
});