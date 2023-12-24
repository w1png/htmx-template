function ClearFormOnSubmit(event, form) {
  if (!event.detail.successful) return;

  form.reset();
}
