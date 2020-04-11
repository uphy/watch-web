<template>
  <v-container fluid>
    <v-list two-line dark>
      <v-list-tile v-for="job in jobs" :key="job.name" @click="openDialog(job.name)">
        <v-list-tile-avatar>
          <v-icon :class="[job.class]">{{ job.icon }}</v-icon>
        </v-list-tile-avatar>

        <v-list-tile-content>
          <v-list-tile-title>{{ job.name }}</v-list-tile-title>
          <v-list-tile-sub-title>{{ job.description }}</v-list-tile-sub-title>
        </v-list-tile-content>

        <v-list-tile-action>
          <v-btn icon ripple>
            <v-icon color="grey lighten-1">info</v-icon>
          </v-btn>
        </v-list-tile-action>
      </v-list-tile>

      <v-divider inset></v-divider>
    </v-list>
    <v-dialog v-model="dialog" max-width="500">
      <v-card>
        <v-card-title class="headline">
          {{ job.name }}
          <v-btn icon ripple @click="checkJob(job.name)">
            <v-icon color="grey lighten-1">autorenew</v-icon>
          </v-btn>
        </v-card-title>

        <v-card-text>
          <v-form>
            <v-container>
              <v-text-field readonly disabled v-model="job.status" label="Status"></v-text-field>
              <v-text-field readonly disabled v-model="job.last" label="Last execution"></v-text-field>
              <v-text-field readonly disabled v-model="job.count" label="Total executions"></v-text-field>
            </v-container>
          </v-form>
          <v-expansion-panel>
            <v-expansion-panel-content>
              <div slot="header">Previous text</div>
              <v-card>
                <v-card-text>{{ job.previous }}</v-card-text>
              </v-card>
            </v-expansion-panel-content>
          </v-expansion-panel>
          <a @click="testJobActions(job.name)">Test alert</a>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="green darken-1" flat="flat" @click="closeDialog()">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import client from "../client";

@Component
export default class App extends Vue {
  jobs: any[] = [];
  dialog: boolean = false;
  job: any = {
    name: "job2",
    description: "Updated on 2018/8/17",
    last: "2018/8/17",
    previous: "aiueo"
  };
  mounted() {
    this.updateList();
  }
  async updateList() {
    let jobs = await client.listJobs();
    jobs = jobs.map(v => {
      switch (v.status) {
        case "ok":
          v.icon = "done";
          v.class = "green--text darken-4";
          break;
        case "running":
          v.icon = "autorenew";
          v.class = "green darken-4";
          break;
        case "error":
          v.icon = "priority_high";
          v.class = "red darken-4";
          break;
      }
      v.description = `Last check ${v.last}, Total ${v.count} times checked.`;
      return v;
    });
    this.jobs = jobs;
  }
  async openDialog(name: string) {
    this.job = await client.getJob(name);
    this.dialog = true;
  }
  closeDialog() {
    this.dialog = false;
  }
  async checkJob(name: string) {
    await client.checkJob(name);
    this.job = await client.getJob(name);
    this.updateList();
  }
  async testJobActions(name:string){
    await client.testJobActions(name);
  }
}
</script>
