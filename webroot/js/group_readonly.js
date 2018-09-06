function httpGroupReadOnly(groupId,readOnly){
    var str="id="+groupId+"&readOnly="+readOnly;
    var result
    var error="null"
    $.ajax({
        url: "/group/readonly",  
        type: "POST",
        data: str,
        async:false,
        dataType:"json", 
        success: function( response ) {
            if(response.result==0){
                result=0
            }else{
                result=response.result
                error=response.info
            }  
        },
        error: function( jqXHR, textStatus, errorThrown ) {
            result=-1
        }
    });

    return [result,error]
}
