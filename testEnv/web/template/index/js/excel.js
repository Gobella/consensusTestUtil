/**
 * Created by gobella on 16-9-23.
 */
/**
 * 将html的table转成Excel的data协议类型数据，不支持ie
 *   table 是HTML DOM Document 对象
 × name 是sheet的名称
 */
var tableToExcel = (function() {
    var uri = 'data:application/vnd.ms-excel;base64,',
        template = '<html xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:x="urn:schemas-microsoft-com:office:excel" xmlns="http://www.w3.org/TR/REC-html40">'
            + '<head><meta http-equiv="Content-type" content="text/html;charset=UTF-8" /><!--[if gte mso 9]><xml><x:ExcelWorkbook><x:ExcelWorksheets><x:ExcelWorksheet><x:Name>{worksheet}</x:Name><x:WorksheetOptions><x:DisplayGridlines/>'
            + '</x:WorksheetOptions></x:ExcelWorksheet></x:ExcelWorksheets></x:ExcelWorkbook></xml><![endif]--></head><body><table>{table}</table></body></html>',
        base64 = function(s) {
            return window.btoa(unescape(encodeURIComponent(s)))
        },
        format = function(s, c) {
            return s.replace(/{(\w+)}/g, function(m, p) {
                return c[p];
            })
        };

    return function(table, name) {
        var ctx = {
            worksheet : name || 'Worksheet',
            table : table.innerHTML
        }
        return uri + base64(format(template, ctx));
    }
})();

$(function(){
    $('#exportExcel').on('click', function(){
        var $this = $(this);
        //设定下载的文件名及后缀
        $this.attr('download', '2016-3-3-财务报表.xls');
        //设定下载内容
        $this.attr('href', tableToExcel($('#targetTable')[0], '财务统计'));
    });
});
