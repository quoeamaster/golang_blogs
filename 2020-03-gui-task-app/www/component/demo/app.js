

new Vue({
    el: '#app',
    data: function () {
        return {
            displayDlg: false,
            notes: []
        };
    },
    methods: {

        // * ------------------ *
        // *   event handlers   *
        // * ------------------ *

        onAddNoteClick: function () {
            this.displayDlg = true;
        },

        onDisplayDlgUpdate: function (data) {
            this.displayDlg = false;
        },
        onAddNote: function (data) {
            this.notes.push(data);
        }

    }
});