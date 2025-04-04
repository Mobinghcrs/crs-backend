// const modal = document.querySelector(".modal");
// const showModal = document.querySelector(".show-modal");
const modalSearch = document.querySelector(".modal-search");
const showModalSearch = document.querySelector(".show-modal-search");
const modalFilter = document.querySelector(".modal-filter");
const showModalFilter = document.querySelector(".show-modal-filter");
const returnModal = document.getElementById("returned-modal");
const closeBtnModal = document.querySelectorAll("#close-modal");

const openModal = (element) => {
  element.classList.remove("hidden");
  element.classList.add("flex");
};

const closeModal = (event, element, innerClass) => {
  if (event.target.closest(innerClass)) return;
  element.classList.add("hidden");
  element.classList.remove("flex");
};
const closedModal = (element) => {
  element.classList.add("hidden");
  element.classList.remove("flex");
};

// showModal.addEventListener("click", () => openModal(modal));
// modal.addEventListener("click", (event) => closeModal(event, modal, ".inner-modal"));

// showModalSearch.addEventListener("click", () => openModal(modalSearch));
// modalSearch.addEventListener("click", (event) =>
//   closeModal(event, modalSearch, ".inner-modal-search")
// );

// showModalFilter.addEventListener("click", () => openModal(modalFilter));
// modalFilter.addEventListener("click", (event) =>
//   closeModal(event, modalFilter, ".inner-modal-filter")
// );

returnModal.addEventListener("click", (event) =>
  closeModal(event, returnModal, ".inner-modal-returned")
);

document.querySelectorAll("#returned").forEach((i) => {
  i.addEventListener("click", () => openModal(returnModal));
});

closeBtnModal.forEach((i) => {
  i.addEventListener("click", () => closedModal(returnModal));
});
