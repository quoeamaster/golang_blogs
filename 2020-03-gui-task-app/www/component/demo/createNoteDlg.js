
Vue.component('create-note-dlg', {
    props: ['show'],
    data: function() {
        return {
            x: "",
            y: "",
            content: ""
        };
    },
    methods: {

        getDlgCss: function () {
            let c = {};

            if (this.show === true) {
                c['core-display-block'] = true;
                c['core-display-none'] = false;
            } else {
                c['core-display-block'] = false;
                c['core-display-none'] = true;
            }
            return c;
        },

        onCancel: function () {
            this.$emit('on-display-dlg-update', false);
            this.resetDialog();
        },
        onSave: function () {
            this.$emit('on-add-note', {
                x: this.x,
                y: this.y,
                content: this.content,
                key: this.content.replace(/\s/g, '')+'_'+(new Date().getMilliseconds())+this.x+'_'+this.y+'_'
            });
            this.$emit('on-display-dlg-update', false);
            this.resetDialog();
        },
        resetDialog: function () {
            this.x = '';
            this.y = '';
            this.content = '';
        }

    },
    template: `
<div style="width: 100%; height: 100%; position: absolute; top: 0; left: 0; z-index: 10; padding-top: 40px; background-color: rgba(230,230,230,0.8);" 
    v-bind:class="getDlgCss()">
    <div style="width: 600px; height: 380px; margin: auto; border: 1px solid #ddd; border-radius: 4px; z-index: 20; overflow: auto; background-color: white;">
        
        <div style="text-align: center; font-size: 1.1em; margin-top: 4px; margin-bottom: 12px;">
            create a note~
        </div>
        
        <div style="padding-left: 12px; padding-right: 12px;">
            <input v-model="x" class="form-control" style="margin-bottom: 4px;" placeholder="position: x coordinate">
            <input v-model="y" class="form-control" style="margin-bottom: 4px;" placeholder="position: y coordinate">
            <textarea v-model="content" class="form-control" rows="5" style="margin-bottom: 4px;" placeholder="content of the note"></textarea>
        </div>
        
        <div class="float-right">
            <button class="btn btn-primary" v-on:click="onSave">save</button>
            <button class="btn btn-secondary" v-on:click="onCancel">cancel</button>
        </div>
    </div>
</div>
    `
});