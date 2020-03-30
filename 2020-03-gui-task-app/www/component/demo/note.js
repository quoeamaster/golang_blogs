
Vue.component('note', {
    props: ['note'],
    mounted: function () {
        let noteObj = document.querySelector("#"+this.note.key);
        if (noteObj) {
            noteObj.style.top = this.note.y+'px';
            noteObj.style.left = this.note.x+'px';
        }
    },
    template: `
<div class='note-container' v-bind:id="this.note.key">
    {{note.content}}
</div>
    `
});