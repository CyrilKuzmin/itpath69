/**
 * PHP Email Form Validation - v3.2
 * URL: https://bootstrapmade.com/changepassword-form/
 * Author: BootstrapMade.com
 */
(function() {
    "use strict";

    let forms = document.querySelectorAll('.changepassword-form');

    forms.forEach(function(e) {
        e.addEventListener('submit', function(event) {
            event.preventDefault();

            let thisForm = this;

            let action = thisForm.getAttribute('action');
            let recaptcha = thisForm.getAttribute('data-recaptcha-site-key');

            if (!action) {
                displayError(thisForm, 'The form action property is not set!')
                return;
            }
            thisForm.querySelector('.loading').classList.add('d-block');
            thisForm.querySelector('.error-message').classList.remove('d-block');
            thisForm.querySelector('.sent-message').classList.remove('d-block');

            let formData = new FormData(thisForm);

            if (recaptcha) {
                if (typeof grecaptcha !== "undefined") {
                    grecaptcha.ready(function() {
                        try {
                            grecaptcha.execute(recaptcha, { action: 'changepassword_form_submit' })
                                .then(token => {
                                    formData.set('recaptcha-response', token);
                                    changepassword_form_submit(thisForm, action, formData);
                                })
                        } catch (error) {
                            displayError(thisForm, error)
                        }
                    });
                } else {
                    displayError(thisForm, 'The reCaptcha javascript API url is not loaded!')
                }
            } else {
                changepassword_form_submit(thisForm, action, formData);
            }
        });
    });

    function changepassword_form_submit(thisForm, action, formData) {
        fetch(action, {
                method: 'POST',
                body: formData,
                headers: { 'X-Requested-With': 'XMLHttpRequest' }
            })
            .then(response => {
                return response.text();
            })
            .then(data => {
                thisForm.querySelector('.loading').classList.remove('d-block');
                if (data.trim() == 'OK') {
                    thisForm.querySelector('.sent-message').classList.add('d-block');
                    thisForm.reset();
                    return;
                } else {
                    throw new Error(data ? data : 'Form submission failed and no error message returned from: ' + action);
                }
            })
            .catch((error) => {
                displayError(thisForm, error);
            });
    }

    function displayError(thisForm, error) {
        thisForm.querySelector('.loading').classList.remove('d-block');
        thisForm.querySelector('.error-message').innerHTML = error.message.split("=")[2].replace('"', '');
        thisForm.querySelector('.error-message').classList.add('d-block');
    }

})();
