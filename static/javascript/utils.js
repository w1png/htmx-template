function ClearFormOnSubmit(event, form) {
  if (!event.detail.successful || event.detail.xhr.status != 200) return;

  form.reset();
}
