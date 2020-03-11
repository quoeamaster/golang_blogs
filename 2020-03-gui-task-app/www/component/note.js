
Vue.component('note', {
    props: ['note'],
    data: function() {
        return {
            'id': ''
        };
    },
    mounted: function() {
        let instance = this;

        // event receiver
        window.eventBus.$on('on-note-dropped', function (data) {
            instance.onNoteDroppedEvent(data);
        });

        this.id = `__note-${this.translateStringToAscii(this.note.content)}-${this.note.x}-${this.note.y}-${this.note.angle}`;
        //console.log("mnt =>", this.id+': ',this.note.content);

        setTimeout(function () {
            let dNote = document.querySelector('#'+instance.id);

            dNote.style.top = instance.note.y + 'px';
            dNote.style.left = instance.note.x + 'px';
            // angle skew
            dNote.style.transform = `rotate(${instance.note.angle}deg)`;

        }, 10);
    },
    methods: {
        /* ---------- */
        /*    util    */
        /* ---------- */

        translateStringToAscii: function (val) {
            let cVal = '';

            for (let i=0; i<val.length; i++) {
                cVal += val.charCodeAt(i)
            }
            return cVal;
        },

        /* ---------------------- */
        /*    event handler(s)    */
        /* ---------------------- */

        onDragStart: function (e) {
            //console.log('dragStart', e);
            window.store.commit('setNoteId', this.id);
            window.store.commit('setOffsetX', e.offsetX);
            window.store.commit('setOffsetY', e.offsetY);
            window.store.commit('setAngle', this.note.angle);
        },

        onNoteDroppedEvent: function (data) {
            if (this.id === data.noteId) {
                console.log('found...', data);
// TODO: update the top, left, angle and call golang API
            }
        }


    },
    template: `
<div class="note-container core-pointer" 
    draggable="true"
    v-on:dragstart="onDragStart(event)"
    v-bind:id="this.id"
    v-bind:key="this.id" >
    <!-- style="top: 0; left: 0;" -->
    <div class="note-inner-container">
        {{note.content}}
    </div>
</div>
`
})