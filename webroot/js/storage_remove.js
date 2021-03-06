function httpStorageRemove(storageId){
    var result
    var error="null"
    var jsonData = {
        "id":storageId
    };
    $.ajax({
        url: "/storage/remove",  
        type: "POST",
        contentType:"application/json",               
        data:JSON.stringify(jsonData), 
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