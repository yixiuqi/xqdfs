var MaskUtil = (function(){  
      
    var $mask,$maskMsg;  
      
    var defMsg = '正在处理，请稍待。。。';  
      
    function init(){  
        if(!$mask){  
            $mask = $("<div></div>")  
            .css({  
              'position' : 'absolute'  
              ,'left' : '0'  
              ,'top' : '0'  
              ,'width' : '100%'  
              ,'height' : '100%'  
              ,'opacity' : '0.3'  
              ,'filter' : 'alpha(opacity=30)'  
              ,'display' : 'none'  
              ,'background-color': '#ECECEC'  
            })  
            .appendTo("body");  
        }  
        if(!$maskMsg){  
            $maskMsg = $("<div></div>")  
                .css({  
                  'position': 'absolute'  
                  ,'top': '50%'  
                  ,'margin-top': '-20px'  
                  ,'padding': '10px 15px 10px 15px'  
                  ,'width': 'auto'  
                  ,'border-width': '2px'  
                  ,'border-style': 'solid'  
				  ,'border-color': '#A8A8A8'  
                  ,'display': 'none'  
                  ,'background-color': '#ECECEC'  
                  ,'font-size':'13px'  
				  ,'color':'#40408B'
                })  
                .appendTo("body");  
        }  
          
        $mask.css({width:"100%",height:$(document).height()});  
          
        var scrollTop = $(document.body).scrollTop();  
          
        $maskMsg.css({  
            left:( $(document.body).outerWidth(true) - 190 ) / 2  
            ,top:( ($(window).height() - 45) / 2 ) + scrollTop  
        });   
                  
    }  
      
    return {  
        mask:function(msg){  
            init();  
            $mask.show();  
            $maskMsg.html(msg||defMsg).show();  
        }  
        ,unmask:function(){  
            $mask.hide();  
            $maskMsg.hide();  
        }  
    }  
      
}());  