<template>
  <div>
    <NodeList />
    <DeploymentList />
    <UnscheduledPodList />
    <PVList />
    <PVCList />
    <StorageClassList />
    <PriorityClassList />
    <NamespaceList />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "@nuxtjs/composition-api";
import NodeList from "./NodeList.vue";
import UnscheduledPodList from "./UnscheduledPodList.vue";
import PVList from "./PVList.vue";
import PVCList from "./PVCList.vue";
import StorageClassList from "./StorageClassList.vue";
import PriorityClassList from "./PriorityClassList.vue";
import NamespaceList from "./NamespaceList.vue";
import DeploymentList from "./DeploymentList.vue";

export default defineComponent({
  components: {
    NodeList,
    UnscheduledPodList,
    DeploymentList,
    PVList,
    PVCList,
    StorageClassList,
    PriorityClassList,
    NamespaceList,
},
});
</script>
