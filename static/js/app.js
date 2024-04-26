
t = 0;
let resp = "";
var converter = new showdown.Converter();

function send(e){
    e.preventDefault();
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


/*async function runScript(prompt, action="/api/v1/chat/stream") {
    id = Math.random().toString(36).substring(2,7);
    outId = "result-" + id;

    $("#printout").append(
        "<div class='px-3 py-3'>" +
        "<div id='" + outId +
        "' style='white-space: pre-wrap;'>" +
        "</div>" +
        "</div>"
    );

    response = await fetch("/api/v1/chat/stream", {
        method: "POST",
        headers: { "Content-Type": "application/json"},
        body: JSON.stringify({content: prompt, llmTimeoutSecond: 30, userId:"1111", maxWindows: 30}),
    });

    decoder = new TextDecoder();
    reader = response.body.getReader();
    while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        $("#"+outId).append(decoder.decode(value));
        window.scrollTo({top: document.body.scrollHeight, behavior:'smooth' });
    }
    $(".js-loading").removeClass("spinner-border");
    $("#"+outId).html(converter.makeHtml($("#"+outId).html()));
    window.scrollTo({top: document.body.scrollHeight, behavior:'smooth' });
    hljs.highlightAll();
}*/

/*async function runScript(prompt) {
    let id = Math.random().toString(36).substring(2, 7);
    let outId = "result-" + id;

    $("#printout").append(
        "<div class='px-3 py-3'>" +
        "<div id='" + outId + "' style='white-space: pre-wrap;'></div>" +
        "</div>"
    );

    try {
        let response = await fetch("/api/v1/chat/stream", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({content: prompt, llmTimeoutSecond: 30, userId:"1111", maxWindows: 30}),
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        let decoder = new TextDecoder();
        let reader = response.body.getReader();
        while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            $("#" + outId).append(decoder.decode(value));
            console.log(decoder.decode(value))
            window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });
        }

        $("#" + outId).html(converter.makeHtml($("#" + outId).html()));
        window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });
        hljs.highlightAll();
    } catch (error) {
        console.error("Fetch error: ", error);
        $("#" + outId).append(`<div>Error loading the response: ${error.message}</div>`);
    } finally {
        $("#sendBtn").attr("disabled", false); // Re-enable the send button
        $(".js-loading").removeClass("spinner-border"); // Assuming you have a spinner with this class
    }
}*/

async function runScript(prompt) {
    let id = Math.random().toString(36).substring(2, 7);
    let outId = "result-" + id;
    let accumulatedContent = "";

    $("#printout").append(
        `<div class='px-3 py-3'><div id='${outId}' style='white-space: pre-wrap;'></div></div>`
    );

    try {
        const response = await fetch("/api/v1/chat/stream", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ content: prompt, llmTimeoutSecond: 30, userId: "1111", maxWindows: 30 }),
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const decoder = new TextDecoder();
        const reader = response.body.getReader();
        do {
            ({ done, value } = await reader.read());
            if (value) {
                let chunkText = decoder.decode(value, { stream: true });
                const json = JSON.parse(chunkText);
                console.log('json:', json)
                accumulatedContent +=  json.data;
                // Update the DOM for every chunk received
                $("#" + outId).html(converter.makeHtml(accumulatedContent));
            }
        } while (!done);

        window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });
        hljs.highlightAll(); // Consider highlighting only after all content is received, or implement more efficient partial highlighting.
    } catch (error) {
        console.error("Fetch error: ", error);
        $("#" + outId).append(`<div>Error loading the response: ${error.message}</div>`);
    } finally {
        $(".js-loading").removeClass("spinner-border");
    }
}



async function clearContent( action) {

    response = await fetch("/api/v1/chat/stream/clear", {
        method: "POST",
        headers: { "Content-Type": "application/json"},
        body: JSON.stringify({content: prompt, llmTimeoutSecond: 30, userId:"1111", maxWindows: 30}),
    });

    console.log(response)

    $("#printout").html("");

    $(".js-logo").removeClass("active");
}
