
Vue.component('create-note-dlg', {
    props: ['show', 'today'],
    data: function() {
        return {
            exitEffects: ['rotateOut', 'rotateOutDownLeft', 'rotateOutDownRight', 'rotateOutUpLeft', 'rotateOutUpRight'],
            exitEffectChosen: '',
            note: {
                content: '',
                x: null,
                y: null,
                angle: null
            }
        };
    },
    methods: {
        /* ------------------------- */
        /*    css [class / style]    */
        /* ------------------------- */

        getBackgroundClass: function () {
            let _c = {};
            if (true === this.show) {
                _c['core-display-block'] = true;
                _c['core-display-none'] = false;
            } else {
                _c['core-display-block'] = false;
                _c['core-display-none'] = true;
            }
            return _c;
        },
        getExitClass: function() {
            let _c = {};
            _c[this.exitEffectChosen] = true;
            return _c;
        },


        /* ----------------------*/
        /*    event handler(s)   */
        /* ----------------------*/

        raiseCloseDlgEvent: function () {
            this.$emit('close-note-creation', false);
        },
        onCreateEvent: function () {
            // create the note / task (etc) => calling Golang backend
            if (!this.note.x) { this.note.x = 0; }
            if (!this.note.y) { this.note.y = 0; }
            if (!this.note.angle) { this.note.angle = 0; }

            window.onCreateNoteTask(this.note.content,
                this.today,
                this.note.x+'', this.note.y+'', this.note.angle+'');

            // add exit effect
            this.exitEffectChosen = this.getRandomExitClass();

            this.resetExitEffectNRaiseEvent();
        },
        onCancelEvent: function () {
            // add exit effect
            this.exitEffectChosen = this.getRandomExitClass();

            this.resetExitEffectNRaiseEvent();
        },
        
        
        /* ---------- */
        /*    util    */
        /* ---------- */
        
        getRandomExitClass: function () {
            let idx = Math.floor(Math.random() * 5);
            return this.exitEffects[idx];
        },
        resetExitEffectNRaiseEvent: function () {
            let instance = this;
            setTimeout(function () {
                // reset the form
                instance.note.content = '';
                instance.note.x = null;
                instance.note.y = null;
                instance.note.angle = null;
                instance.raiseCloseDlgEvent();
                instance.exitEffectChosen = ''; // reset
            }, 1000);
        }

    },
    template: `
<div class="note-cr-dlg-background" v-bind:class="getBackgroundClass()">
    <div class="note-cr-dlg-container animated" v-bind:class="getExitClass()">
        <!-- header -->
        <div class="note-cr-dlg-header">
            <span class="float-right note-cr-close-btn core-pointer" v-on:click="onCancelEvent">&times;</span>
            <span class="float-left">
                create a task note~
            </span>    
        </div>
        
        <!-- main content -->
        <div class="note-cr-content-outer-container">
            <div class="note-cr-content-inner-container">
                <div class="note-cr-flex3-row">
                    <input type="text" class="form-control note-cr-flex3-text" v-model="note.x" 
                        style="width: auto;"  placeholder="x / left position" />
                    <input type="text" class="form-control note-cr-flex3-text" v-model="note.y" 
                        style="width: auto;"  placeholder="y / top position" />
                    <input type="text" class="form-control note-cr-flex3-text" v-model="note.angle" 
                        style="width: auto;"  placeholder="angle" />
                </div>
                
                <textarea placeholder="task / note content..."
                          style="margin-top: 8px;" 
                          v-model="note.content"
                          class="form-control note-cr-content" rows="4"></textarea>
            </div>
        </div>
        
        <!-- footer -->
        <button v-on:click="onCreateEvent" class="btn btn-primary">create</button>
        <button v-on:click="onCancelEvent" class="btn btn-secondary">nah... not now</button>
    </div>
</div>
    `
});