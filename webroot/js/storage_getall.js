function httpStorageGetAll(){
    var result
    var error="null"
    var data
    $.ajax({
        url: "/storage/getall",  
        type: "POST",
        data: "",
        async:false, 
        dataType:"json", 
        success: function( response ) {
            if(response.result==0){
                result=0
                data=response.rows
            }else{
                result=response.result
                error=response.info
            }            
        },
        error: function( jqXHR, textStatus, errorThrown ) {
            result=-1
        }
    });

    return [result,error,data]
}