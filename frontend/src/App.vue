<template>
  <v-app>
    <v-toolbar app clipped-left>
      <v-toolbar-title>watch bot console</v-toolbar-title>
    </v-toolbar>
    <v-content>
      <v-card>
        <v-container fluid grid-list-lg>
          <v-layout row wrap>
            <v-flex xs12>
              <v-card>
                <v-card-title primary-title>
                  <div class="headline">
                    Jobs
                    <v-btn icon ripple @click="reloadJobs()">
                      <v-icon>autorenew</v-icon>
                    </v-btn>
                  </div>
                  <Jobs ref="jobs" />
                </v-card-title>
              </v-card>
            </v-flex>
          </v-layout>
        </v-container>
      </v-card>
    </v-content>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Jobs from "./components/Jobs.vue";

@Component({
  components: {
    Jobs
  }
})
export default class App extends Vue {
  mounted() {}
  reloadJobs() {
    let jobs:any = this.$refs.jobs;
    jobs.updateList()
  }
}
</script>
