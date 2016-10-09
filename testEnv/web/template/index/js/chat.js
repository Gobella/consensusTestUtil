/**
 * Created by gobella on 16-9-27.
 */
$(function () {
    $('#flush').click(function () {

        $.ajax({
            url: "http://127.0.0.1:8880/statistic/count",
            context: $('#targetTable'),
            type: "post",
            dataType: "json",
            data:{"type":"count","id":2},
            success: function(data){
                var jsonText = JSON.stringify(data);
                var jsonObj=JSON.parse(jsonText)
                alert(jsonObj.Nodeinfo.Id)
            }});
    });

});

