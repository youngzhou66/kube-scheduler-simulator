import { V1Deployment } from "@kubernetes/client-node";
import { reactive, inject } from "@nuxtjs/composition-api";
import { DeploymentAPIKey } from "~/api/APIProviderKeys";

type stateType = {
  selectedDeployment: selectedDeployment | null;
  deployments: V1Deployment[];
  lastResourceVersion: string;
};

export type selectedDeployment = {
  // isNew represents whether this is a new one or not.
  isNew: boolean;
  item: V1Deployment;
  resourceKind: string;
  isDeletable: boolean;
};


export default function deploymentStore() {
  const state: stateType = reactive({
    selectedDeployment: null,
    deployments: [],
    lastResourceVersion: "",
  });
  const deploymentAPI = inject(DeploymentAPIKey);
  if (!deploymentAPI) {
    throw new Error(`${DeploymentAPIKey.description} is not provided`);
  }
  return {
    get count(): number {
      return state.deployments.length;
    },

    get selected() {
      return state.selectedDeployment;
    },

    select(d: V1Deployment | null, isNew: boolean) {
      if (d !== null) {
        state.selectedDeployment = {
          isNew: isNew,
          item: d,
          resourceKind: "Deployment",
          isDeletable: true,
        };
      }
    },

    // initList calls list API, and stores current resource data and lastResourceVersion.
    async initList() {
      const list = await deploymentAPI.listDeployment();
    },

    resetSelected() { },
    async apply(d: V1Deployment) {
      console.log(123123)
      await deploymentAPI.createDeployment(d)
    },
    async fetchSelected() { },
    async delete(d: V1Deployment) { },

  };
}

export type DeploymentStore = ReturnType<typeof deploymentStore>;