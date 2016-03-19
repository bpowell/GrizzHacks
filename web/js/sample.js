function poop() {
    $.ajax({
        type: "POST",
    url: 'http://104.131.18.185:8080/api/getday',
        data: { ticker: "googl", date: "2016-Mar-18" },
    success: success,
    dataType: 'json'
    });
}

function success(a, b, c) {
    console.log(a);
    console.log(b);
    console.log(c);
}
