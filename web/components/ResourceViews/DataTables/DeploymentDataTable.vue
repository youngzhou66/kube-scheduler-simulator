<template>
  <DataTable
      :title="`Deployments`"
      :headers="headers"
      :items="deployments"
      :on-click="onClick"
  />
</template>

<script lang="ts">
import { V1Deployment } from "@kubernetes/client-node";
import { computed, inject, defineComponent } from "@nuxtjs/composition-api";
import DataTable from "./DataTable.vue";
import {} from "../../lib/util";
import DeploymentStoreKey from "~/components/StoreKey/DeploymentStoreKey";

export default defineComponent({
  components: {
    DataTable,
  },
  setup() {
    const store = inject(DeploymentStoreKey);
    if (!store) {
      throw new Error(`${DeploymentStoreKey.description} is not provided`);
    }

    const onClick = (depolyment: V1Deployment) => {
      store.select(depolyment,false);
    };

    const pods = computed(() => {
      return {};
    });
    const search = "";
    const headers = [
      {
        text: "Name",
        value: "metadata.name",
        sortable: true,
      },
      { text: "Namespace", value: "metadata.namespace", sortable: true },
      { text: "Node", value: "spec.nodeName", sortable: true },
      {
        text: "Conditions",
        value: "status.conditions[0].type",
        sortable: true,
      },
      {
        text: "CreationTime",
        value: "metadata.creationTimestamp",
        sortable: true,
      },
      {
        text: "UpdateTime",
        value: "metadata.managedFields[0].time",
        sortable: true,
      },
    ];
    return {
      pods,
      search,
      onClick,
      headers,
    };
  },
});
</script>