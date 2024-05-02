
t = 0;
let resp = "";
var converter = new showdown.Converter();

function send(e){
    e.preventDefault();
    chatStream()
}

$(document).ready(function(){
    $('#prompt').keypress(function(event){
        var keycode = (event.keyCode ? event.keyCode : event.which);
        if((keycode == 10 || keycode == 13) && event.ctrlKey){
            send(event);
            return false;
        }
    });
    autosize($('#prompt'));
});

$(document).ready(function(){
    $('.clear-button').click(function() {
        clearContent();
    });
});

$(document).ready(function () {
    $('.send-button').click(function () {
        chatStream()
    })
})

function chatStream() {
    var prompt = $("#prompt").val().trimEnd();
    $("#prompt").val("");
    autosize.update($("#prompt"));

    $("#printout").append(
        "<div class='prompt-message'>" +
        "<div style='white-space: pre-wrap;'>" +
        prompt  +
        "</div>" +
        "<span class='message-loader js-loading spinner-border'></span>" +
        "</div>"
    );
    window.scrollTo({top: document.body.scrollHeight, behavior:'smooth' });
    runScript(prompt);
    $(".js-logo").addClass("active");
}

async function runScript(prompt) {
    let id = Math.random().toString(36).substring(2, 7);
    let outId = "result-" + id;
    let accumulatedContent = "";

    disableButtons();

    if (prompt.trim() === "") {
        return
    }

    $("#printout").append(
        `<div class='px-3 py-3'><div id='${outId}' style='white-space: pre-wrap;'></div></div>`
    );

    try {
        const response = await fetch("/api/v1/chat/stream", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ content: prompt, llmTimeoutSecond: 30, userId: "1111", maxWindows: 30, llmModel: "qwen:14b" }),
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const decoder = new TextDecoder()
        const reader = response.body.getReader()
        do {
            ({ done, value } = await reader.read())
            if (value) {
                let chunkText = decoder.decode(value, { stream: true })
                console.log(chunkText)
                var lines = chunkText.split('\n').filter(line => line !== '')
                for (const line of lines){
                    const json = JSON.parse(line);
                    if (json.code === 'S00000') {
                        accumulatedContent +=  json.data;
                        $("#" + outId).html(converter.makeHtml(accumulatedContent));
                    } else {
                        accumulatedContent += json.message
                        $("#" + outId).html(converter.makeHtml(accumulatedContent));
                        break
                    }
                }
            }
        } while (!done);

        window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });
        hljs.highlightAll(); // Consider highlighting only after all content is received, or implement more efficient partial highlighting.
    } catch (error) {
        console.error("Fetch error: ", error);
        $("#" + outId).append(`<div>Error loading the response: ${error.message}</div>`);
    } finally {
        $(".js-loading").removeClass("spinner-border");
        enableButtons()
    }
}

function disableButtons() {
    $('.send-button, .clear-button').prop('disabled', true);
}

function enableButtons() {
    $('.send-button, .clear-button').prop('disabled', false);
}


async function clearContent( action) {
    // 禁用按钮
    disableButtons();

    response = await fetch("/api/v1/chat/stream/clear", {
        method: "POST",
        headers: { "Content-Type": "application/json"},
        body: JSON.stringify({content: prompt, llmTimeoutSecond: 30, userId:"1111", maxWindows: 30}),
    });

    console.log(response)

    $("#printout").html("");

    $(".js-logo").removeClass("active");

    // 操作完成，启用按钮
    enableButtons();
}
