import { Modal } from "./modal";

function register(tag, component) {
  customElements.define(tag, component);
}

export {
  register,

  /* components */
  Modal,
};
