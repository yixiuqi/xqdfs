function httpStorageInit(storageId){
    var str="id="+storageId;
    var result
    var error="null"
    $.ajax({
        url: "/storage/init",  
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