
Vue.component('note', {
    props: ['note'],
    data: function() {
        return {
            'id': ''
        };
    },
    mounted: function() {
        this.id = `__note-${this.translateStringToAscii(this.note.content)}-${this.note.x}-${this.note.y}-${this.note.angle}`;
        //console.log("mnt =>", this.id+': ',this.note.content);

        let instance = this;
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
        }



    },
    template: `
<div class="note-container core-pointer" 
    v-bind:id="this.id"
    v-bind:key="this.id" >
    <!-- style="top: 0; left: 0;" -->
    <div class="note-inner-container">
        {{note.content}}
    </div>
</div>
`
})