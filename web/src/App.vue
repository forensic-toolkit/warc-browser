<template>
  <w-app>

    <!-- header -->
    <w-toolbar color="grey-dark4">
      <div class="title2 text-italic">warc-browser</div>
    </w-toolbar>

    <!-- body -->
    <w-flex row>

      <!-- left-side file listing -->
      <div class="xs3">
        <div class="py3">
          <w-tree
            deep-reactivity
            v-model="selection"
            :data="tree"
            @click="treeclick"
            expand-icon="mdi mdi-plus"
            expand-open-icon="mdi mdi-minus"
            branch-icon="mdi mdi-folder-outline"
            branch-open-icon="mdi mdi-folder-open-outline"
            class="grey-dark1 title">
          </w-tree>
        </div>
      </div>

      <!-- right-side warc viewer -->
      <div class="xs11" fill-height>
          <w-transition-expand y>
            <div class="transition-box fill-height" v-if="selected">
              <w-card color="grey-dark1" class="fill-height">
                <replay-web-page v-if="warc_url" :source="warc_url" class="fill-height fill-width">
                </replay-web-page>
                <template #actions>
                  <div class="spacer"></div>                 
                  <w-button text outline :route="warc_url">
                    {{ warc_label }}
                    <w-icon class="ml2" sm>mdi mdi-open-in-new</w-icon>
                  </w-button>
                </template>
              </w-card>
            </div>
          </w-transition-expand>
      </div>
    </w-flex>
  </w-app>
</template>

<script>
export default {
  mounted() {
    this.fetchdir("archives/").then((data)=>{
      if ( data.error == null ) {
        this.tree = data.records;
      }
    });
  },
  methods: {
    async fetchdir(path) {
      console.log(`[fetchdir][${path}]`);
      const resp = await fetch(`/${path}`)
      const data = await resp.json()
      console.log(`[fetchdir][${path}]:`, data);
      return data
    },
    treeclick(e) {      
      console.log('[treeclick].selection', e.item);
      if ( e.item.branch ) {
        this.fetchdir(e.item.path).then((data)=>{
          if ( data.error == null ) {
            e.item.children = data.records;
          }
        })
      } else {
        this.selected = false;
        this.warc_url = e.item.path;
        this.warc_label = e.item.warc_label;
        this.$nextTick().then(()=>{
          this.selected = true;
        })
      }
    },
  },
  data: () => ({
    tree: [],
    selected: null,
    selection: null,
    warc_url: "",
    warc_label: "",    
  })
}
</script>