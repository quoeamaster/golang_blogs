
Vue.component('splash-screen', {
    data: function() {
        return {
            show: true
        };
    },
    mounted: function() {
        let instance = this;
        /*
        setTimeout(function () {
            instance.show = false;
        }, 2000);
        */

        let ePBar = document.querySelector('#splash-progress-bar');
        let w = 1;
        let hInterval = setInterval(function () {
            if (w >= 100) {
                clearInterval(hInterval);
                instance.show = false;
            } else {
                ePBar.style.width = w + '%';
            }
            w++;
        }, 50);

    },
    methods: {
        getContainerClass: function () {
            let _c = {};
            if (this.show) {
                _c['core-display-block'] = true;
                _c['core-display-none'] = false;
            } else {
                _c['core-display-block'] = false;
                _c['core-display-none'] = true;
            }
            return _c;
        }
    },
    template: `
<div class="splash-container" v-bind:class="getContainerClass()">
    <div class="splash-inner-container">
        <div>
            <div class="splash-caption" 
                style="margin-top: 180px; margin-bottom: 20px;">
                <img src="../asset/favicon.png" width="50px; margin-right: 12px;">
                <span>Loading...</span>
            </div>
            
            <div class="splash-progress-container">
                <div id='splash-progress-bar' class="splash-progress-bar"></div>
            </div>
        </div>
    </div>
</div>
    `
});