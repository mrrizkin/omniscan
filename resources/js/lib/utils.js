/**
 * Registers a callback function to be executed when the DOM content is fully loaded.
 *
 * @param {Function} callback - The function to execute once the DOM is ready.
 */
export function onMounted(callback) {
    document.addEventListener("DOMContentLoaded", callback);
}
