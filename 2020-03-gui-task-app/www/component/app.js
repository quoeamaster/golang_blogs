// create a global event-bus
let eventBus = new Vue();
window.eventBus = eventBus;

// create a vuex store
window.store = new Vuex.Store({
    state: {
        // the draggable note's key/id, __none__ means no note dragging or dragged before
        noteId: '__none__',

        // actual offset(s) on dragging the note (offset = distance between the top,left of a note)
        offsetX: 0,
        offsetY: 0,
        angle: 0
    },
    mutations: {
        setNoteId(state, id) {
            state.noteId = id;
        },
        setOffsetX(state, x) {
            state.offsetX = x;
        },
        setOffsetY(state, y) {
            state.offsetY = y;
        },
        setAngle(state, angle) {
            state.angle = angle;
        }
    }
});



new Vue({
    el: '#app',
    data: function () {
        return {
            showSideMenu: false,

            today: new Date(),
            todayInString: '',
            chosenDateToDisplay: '',

            notes: {},
            //chosenNotes: {},

// TODO: are there any contextMenu showing
            isContextMenuShowing: false,
        };
    },
    mounted: function() {
        // get all notes from repo
        window.onGetNotes();

        // callback on getting notes OR after notes being updated (update the "notes" object)
        window.eventBus.$on('go-notes-all-loaded', this.onGoNotesAllLoadedEvent);

        if (!this.today) {
            // should not happen though
            this.today = new Date();
        }
        this.todayInString = this.today.getFullYear() + '-';
        this.todayInString += ((this.today.getMonth()+1)<10)?'0'+(this.today.getMonth()+1):this.today.getMonth()+1;
        this.todayInString += "-";
        this.todayInString += (this.today.getDate()<10)?'0'+this.today.getDate():this.today.getDate();
        // for JUST mounted -> today = chosen-date
        this.chosenDateToDisplay = this.todayInString;
        //this.chosenNotes = this.notes[this.chosenDateToDisplay];


// TODO: add back a click event to remove our custom contextMenu
        window.addEventListener('click', function (event) {
            if (this.isContextMenuShowing === true) {
// TODO: HIDE it
            }
        });

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
            this.isContextMenuShowing = true;
        },
        onCloseNoteCreation: function () {
            this.isContextMenuShowing = false;
        },
        onDateChangeByDelta: function (delta) {
            // chosenDateToDisplay => convert back to date object
            // add the above date with the delta
            // translate the above new date back to string and update the "chosenDateToDisplay" variable
            // load the corresponding notes as well... (usually no need to trigger another golang operation since the notes are already in memory)
            let dParts = this.chosenDateToDisplay.split("-");
            if (dParts.length === 3) {
                let month = parseInt(dParts[1], 10) - 1;
                let dChosen = new Date(parseInt(dParts[0], 10), month, parseInt(dParts[2], 10), 0,0,0);

                dChosen.setDate(dChosen.getDate() + delta);
                let dChosenInString = dChosen.getFullYear() + '-';
                dChosenInString += ((dChosen.getMonth()+1)<10)?"0"+(dChosen.getMonth()+1):dChosen.getMonth()+1;
                dChosenInString += "-";
                dChosenInString += ((dChosen.getDate())<10)?"0"+(dChosen.getDate()):dChosen.getDate();

                this.chosenDateToDisplay = dChosenInString;
                //this.chosenNotes = this.notes[this.chosenDateToDisplay];
            }

        },


        /* ------------------------------ */
        /*    golang interfacing event    */
        /* ------------------------------ */

        onGoNotesAllLoadedEvent: function (data) {
            this.notes = data;
            // TODO: debug =>
            // console.log('** inside notesLoaded', data);
        }


    },
});