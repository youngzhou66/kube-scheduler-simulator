import { InjectionKey } from "@nuxtjs/composition-api";
import { DeploymentStore } from "../../store/deployment";

const DeploymentStoreKey: InjectionKey<DeploymentStore> = Symbol("DeploymentStore");
export default DeploymentStoreKey;
