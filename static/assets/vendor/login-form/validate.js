(function() {
    "use strict";

    let forms = document.querySelectorAll('.login-form');

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
                            grecaptcha.execute(recaptcha, { action: 'login_form_submit' })
                                .then(token => {
                                    formData.set('recaptcha-response', token);
                                    login_form_submit(thisForm, action, formData);
                                })
                        } catch (error) {
                            displayError(thisForm, error)
                        }
                    });
                } else {
                    displayError(thisForm, 'The reCaptcha javascript API url is not loaded!')
                }
            } else {
                login_form_submit(thisForm, action, formData);
            }
        });
    });

    function login_form_submit(thisForm, action, formData) {
        fetch(action, {
                method: 'POST',
                body: formData,
                headers: { 'X-Requested-With': 'XMLHttpRequest' }
            })
            .then(function(response) {
                if (response.ok) {
                    window.location.href = "/lk";
                } else {
                    return response.json();
                }
            })
            .then(function(json) {
                if (json != null) {
                    throw new Error(`${json.message}`);
                }
            })
            .catch(function(error) {
                displayError(thisForm, error);
            });
    }

    function displayError(thisForm, error) {
        thisForm.querySelector('.loading').classList.remove('d-block');
        thisForm.querySelector('.error-message').innerHTML = error;
        thisForm.querySelector('.error-message').classList.add('d-block');
    }

})();