
Vue.component('note', {
    props: ['note'],
    data: function() {
        return {
            'id': ''
        };
    },
    mounted: function() {
        this.id = `__note-${this.translateStringToAscii(this.note.content)}-${this.note.x}-${this.note.y}-${this.note.angle}`;

        let instance = this;
        setTimeout(function () {
            let dNote = document.querySelector('#'+instance.id);

            dNote.style.top = instance.note.y + 'px';
            dNote.style.left = instance.note.x + 'px';
            // angle skew
            dNote.style.transform = `rotate(${instance.note.angle}deg)`;

        }, 100);
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
        }


        /*
        getNoteLocationStyle: function () {
            let _s = {};
            _s['top'] = this.note.y + 'px;';
            _s['left'] = this.note.x + 'px;';
            // TODO: skew angle...

            return JSON.stringify(_s);
        }*/

    },
    template: `
<div class="note-container core-pointer" v-bind:id="this.id">
    <div class="note-inner-container">
        {{note.content}}
    </div>
</div>
`
})