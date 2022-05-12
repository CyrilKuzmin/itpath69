/*
Как оно должно работать.
При создании нового теста (заходим на страницу testing, жмем кнопку "Начать"):
* Запрашиваем JSON с тестом
* Генерируем слайды по вопросам и слайд-бар, заполняем карусель
* Сохраняем тест в localStorage по имени test-{{ ID }}-body
* Вешаем на чекбоксы и радио ивенты, которые будут сохранять текущий ответ в localStorage 
  по имени test-{{ ID }}-ans-{{ # of Q }}. Ответ содержит список выбранных text'ов (на случай если поряок будет изменен)
* При нажатии "Проверки" делаем POST /test, передавая сформированный body с текущими ответами (is_correct ставим в true на выбранных)
* Скрываем карусель и показываем результат (ну или ошибку, да)
* Чистим localStorage

Если открываем уже имеющийся в localStorage тест, то не запрашиваем JSON, а просто стаивм уже выбранные ответы, если есть
*/

var AnswersCounter = 0
var TotalQuestions = 0

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function toogleCarouselVisibility(makeVisible) {
    if (makeVisible) {
        document.getElementById("btn-start").remove() // we'll not shwo this BTN again
        document.getElementById("testing-logo").classList.add("d-none");
        document.getElementById("testing-intro").classList.add("d-none");
        document.getElementById("testing-check-btn").classList.remove("d-none");
        document.getElementById("testing-answers-counter").classList.remove("d-none");
        document.getElementById("testing").classList.remove("d-none");
    } else {
        document.getElementById("testing-logo").classList.remove("d-none");
        document.getElementById("testing-intro").classList.remove("d-none");
        document.getElementById("testing-check-btn").classList.add("d-none");
        document.getElementById("testing-answers-counter").classList.add("d-none");
        document.getElementById("testing").classList.add("d-none");
    }
}

function displayScore(score) {
    result_elem = document.getElementById("score-result")
    details_elem = document.getElementById("score-details")
    details_elem.innerHTML = `Вы набрали ${(score*100).toFixed(2)}%`
    if (score > 0.85) {
        result_elem.innerHTML = `Поздравляю! Тест пройден`
    } else {
        result_elem.innerHTML = `Вы не прошли тест`
    }
}

function updateCounter() {
    document.getElementById("testing-answers-counter").innerHTML = `${AnswersCounter}/${TotalQuestions}`
    btn = document.getElementById("testing-check-btn")
    if (AnswersCounter == TotalQuestions) {
        btn.classList.remove("btn-warning")
        btn.classList.add("btn-success")
    } else {
        btn.classList.remove("btn-success")
        btn.classList.add("btn-warning")
    }
}

function generateCarouselItem(idx, question) {
    res = `
    <div id="q${idx}" class="carousel-item">
        <div class="carousel-caption">
            <div class="container">
                ${generateQuestionImage(idx, question)}
                <h5>${question.question_text}</h5>
                ${generateAnswerOptions(idx, question)}
            </div>
        </div>
    </div>
    `
    return res
}

function generateQuestionImage(idx, question) {
    if (question.image_url == "") {
        return ""
    }
    return `
    <div class="module-img">
        <a href="${question.image_url}" data-lightbox="image-${idx}"><img src="${question.image_url}" class="img-fluid"></a>
    </div>
    `
}

function generateAnswerOptions(idx, question) {
    res = ""
    for (var i = 0; i < question.answers.length; i++) {
        switch (question.question_type) {
            case 1:
                res += generateRadio(idx, i, question.answers[i].text)
                break
            case 2:
                res += generateCheckbox(idx, i, question.answers[i].text)
                break
        }
    }
    return res
}

function generateCheckbox(q_id, a_id, text) {
    return `
    <div class="form-check">
        <input class="form-check-input answer-checkbox" type="checkbox" id="q-${q_id}-a-${a_id}">
        <label class="form-check-label" for="q-${q_id}-a-${a_id}" id="q-${q_id}-a-${a_id}-label">${text}</label>
    </div>
    `
}

function generateRadio(q_id, a_id, text) {
    return `
    <div class="form-check">
        <input class="form-check-input answer-radio" type="radio" name="${q_id}-radios" id="q-${q_id}-a-${a_id}" >
        <label class="form-check-label" for="q-${q_id}-a-${a_id}" id="q-${q_id}-a-${a_id}-label">${text}</label>
    </div>
    `
}

function generateSlideIndicator(q_id) {
    return `
    <button id="q${q_id}-indic" type="button" data-bs-target="#testing-carousel" data-bs-slide-to="${q_id}" 
    aria-label="Question ${q_id}" aria-current="true"></button>
    `
}

function checkboxChanged(event, test_id, item_id, store) {
    q_id = item_id.split("-").at(1);
    label_id = item_id + '-label'
    checkbox = document.getElementById(item_id)
    text = document.getElementById(label_id).innerHTML
    storageItem = `test-${test_id}-ans-${q_id}`
    if (checkbox.checked) {
        answers = store.getItem(storageItem)
        if (answers == null) {
            AnswersCounter++
            store.setItem(storageItem, JSON.stringify([text]))
            updateCounter()
        } else {
            myAnswers = JSON.parse(answers)
            if (myAnswers.length == 0) {
                AnswersCounter++
                updateCounter()
            }
            myAnswers.push(text)
            store.setItem(storageItem, JSON.stringify(myAnswers))
        }
    } else {
        answers = store.getItem(storageItem)
        if (answers == null) {
            return
        } else {
            myAnswers = JSON.parse(answers)
            if (myAnswers.length == 1) {
                AnswersCounter--
                updateCounter()
            }
            index = myAnswers.indexOf(text);
            if (index > -1) {
                myAnswers.splice(index, 1); // 2nd parameter means remove one item only
                store.setItem(storageItem, JSON.stringify(myAnswers))
            }
        }
    }
}

function radioClicked(event, test_id, item_id, store) {
    q_id = item_id.split("-").at(1);
    label_id = item_id + '-label'
    text = document.getElementById(label_id).innerHTML
    storageItem = `test-${test_id}-ans-${q_id}`
    currentAnswer = store.getItem(storageItem)
    if (currentAnswer == null) {
        AnswersCounter++
        updateCounter()
    }
    store.setItem(storageItem, JSON.stringify([text]))
}

function checkResults(event, test_id, store) {
    event.preventDefault();
    if (AnswersCounter != TotalQuestions) {
        var result = confirm("Даны ответы не на все вопросы. Вы уверены, что хотите завершить тест?");
        if (!result) {
            return
        }
    }
    storageTest = `test-${test_id}-body`
    test_body = store.getItem(storageTest)
    if (test_body == null || test_body == "") {
        console.log("test is missing. WTF?")
        return
    }
    test = JSON.parse(test_body)
    for (var q_id = 0; q_id < test.questions.length; q_id++) {
        storageAnswer = `test-${test_id}-ans-${q_id}`
        raw_answer = store.getItem(storageAnswer)
        if (raw_answer == null || raw_answer == "") {
            continue
        }
        chosen = JSON.parse(raw_answer)
        for (var a_id = 0; a_id < test.questions[q_id].answers.length; a_id++) {
            idx = chosen.indexOf(test.questions[q_id].answers[a_id].text)
            if (idx > -1) {
                test.questions[q_id].answers[a_id].is_correct = true
            } else {
                test.questions[q_id].answers[a_id].is_correct = false
            }
        }
        store.removeItem(storageAnswer)
    }
    store.setItem(storageTest, JSON.stringify(test))
    fetch('/test', {
            method: 'POST', // or 'PUT' ?
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(test),
        })
        .then(response => response.json())
        .then(data => {
            toogleCarouselVisibility(false)
            displayScore(data.score)
            store.removeItem(storageTest)
        })
        .catch((error) => {
            console.error('Error:', error);
        });
}

function restore_answers_from_storage(store, test_id) {
    storageTest = `test-${test_id}-body`
    test_body = store.getItem(storageTest)
    if (test_body == null || test_body == "") {
        console.log("test is missing. WTF?")
        return
    }
    test = JSON.parse(test_body)
    for (var q_id = 0; q_id < test.questions.length; q_id++) {
        storageAnswer = `test-${test_id}-ans-${q_id}`
        raw_answer = store.getItem(storageAnswer)
        if (raw_answer == null || raw_answer == "") {
            continue
        }
        chosen = JSON.parse(raw_answer)
        if (chosen.length > 0) {
            AnswersCounter++
        }
        for (var a_id = 0; a_id < test.questions[q_id].answers.length; a_id++) {
            label_text = document.getElementById(`q-${q_id}-a-${a_id}-label`).innerHTML
            if (chosen.indexOf(label_text) > -1) {
                document.getElementById(`q-${q_id}-a-${a_id}`).checked = true;
            }
        }
    }
    updateCounter()
}

function start(store, test_id, module_id) {
    toogleCarouselVisibility(true);
    slides = ""
    slide_indicators = ""
    $.getJSON(`test?module_id=${module_id}&test_id=${test_id}`, function(data) {
        test_id = data.id
            // Generating slides and bar HTML based on JSON received
        TotalQuestions = data.questions.length
        for (var i = 0; i < data.questions.length; i++) {
            slides += generateCarouselItem(i, data.questions[i])
            slide_indicators += generateSlideIndicator(i)
        }
        document.getElementById("testing-inner").innerHTML += slides;
        document.getElementById("slides-indicators").innerHTML += slide_indicators;
        // mark 1st question as active
        document.getElementById("q0").classList.add("active");
        document.getElementById("q0-indic").classList.add("active");
        // save test to storage
        store.setItem(`test-${test_id}-body`, JSON.stringify(data))
            // event listeneres on checkboxes and radios
        document.querySelectorAll('.answer-checkbox').forEach(item => {
            item.addEventListener('change', event => checkboxChanged(event, test_id, item.id, store))
        })
        document.querySelectorAll('.answer-radio').forEach(item => {
            item.addEventListener('click', event => radioClicked(event, test_id, item.id, store))
        })
        updateCounter()
            // event listener for check button
        document.getElementById('testing-check-btn').addEventListener('click', event => checkResults(event, test_id, store))
            // maybe we already have some answers (restart case)
        restore_answers_from_storage(store, test_id)
    });
}

function new_start(module_id) {
    testsStorage = window.localStorage;
    start(testsStorage, "", module_id)
}

function restart(test_id, module_id) {
    testsStorage = window.localStorage;
    start(testsStorage, test_id, module_id)
}