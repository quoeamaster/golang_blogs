
new Vue({
    el: '#app',
    data: function () {
        return {
            showSideMenu: false,

            today: new Date(),
            todayInString: '',
        };
    },
    mounted: function() {
        if (!this.today) {
            // should not happen though
            this.today = new Date();
        }
        this.todayInString = this.today.getFullYear() + '-';
        this.todayInString += ((this.today.getMonth()+1)<10)?'0'+(this.today.getMonth()+1):this.today.getMonth()+1;
        this.todayInString += "-";
        this.todayInString += (this.today.getDate()<10)?'0'+this.today.getDate():this.today.getDate();
    },
    methods: {
        /* ------------------ */
        /*   event handlers   */
        /* ------------------ */

        onMenuClick: function () {
            this.showSideMenu = true;
        },
        onUpdateShowSideMenu: function (data) {
            this.showSideMenu = data;
        },
        onUpsertMemo: function () {
            console.log('on-upsert');
        }


    },
});