$.fn.exists = function () {
    return this.length !== 0;
};

$(document).ready(function () {

    //Toggle extra info on notice class of notification and dismiss



    $('li.notice').click(function () {
        $(this).children('div').slideToggle();
		//get the list of participants
		var x = $(this).children('a.whos-going').attr('data-eventId');
		var y = $(this);
		if (y.children('div.more-info').children('ul').children('li').length === 0) {
			$.ajax('/participants/'+ x)
			.done (function(data){
				for(var i in data) {
					$('<li>'+data[i].name+'</li>').appendTo(y.children('div.more-info').children('ul'));
				}
			})
		}
			
	});
	
	
    $('li.alert-box a.dismiss-notification').click(function (e) {
        e.preventDefault;
        $(this).parent().fadeOut(function () {
            $(this).remove();
        });
    });


    //Add ME to list on click of the join button


    $('input.joinevent-button').on('click', function (e) {
		e.preventDefault();
        var x = $(this).parent().parent();
		var y = x.children('div.more-info').children('ul');
		if (y.children('li.me').length === 0) {
			var event = x.children('a.whos-going').attr('data-eventId');
			$.ajax({
				type: "POST",
				url: '/join',
				data: {id: event},
				success: function(){
					$('<li class="me">ME</li>').appendTo(y);
				}
				
			});		
        }
    });



    //Change Password toggle

    $('a#clickToEditPw').click(function () {
        $('.changePw').slideToggle();
    });


    $('div.alert-box a').click(function (e) {
        e.preventDefault;
        $(this).parent().fadeOut(function () {
            $(this).remove();
        });
    });


    //Change Name

    $('#clickToEditName').on('click', function () {
        if ($(this).html() == "Save") {
            var name = $('#name').val();
            $('#name').replaceWith('<span id="name">' + name + '</span>');
            $(this).html("Edit");
        } else {
            var name = $('#name').html();
            $('#name').replaceWith('<input id="name" value="' + name + '">');
            $(this).html("Save");
        }
    });

	

	

    //Toggle Extra Info


    $('a.whatIsThis').click(function () {
        $('#popupInfo').slideToggle();
    });




    // good emails : green, bad emails : red

    var tag = $('input#emailsUserInput');
    if (tag.exists()) {
        tag.tagsinput({
            tagClass: function (item) {
                if (/\b[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}\b/i.test(item)) {
                    return 'valid';
                } else {
                    return 'invalid';
                }
            }
        });
    }

    // typeahead for circles input

	
	var substringMatcher = function() {
	var circles = null;
	if (!circles) {
		$.get('/circles', {async: false}, function(data) {
			circles = data;
		})
	}
	return function findMatches(q, cb) {
    var matches, substrRegex;
 
		// an array that will be populated with substring matches
		matches = [];
 
		// regex used to determine if a string contains the substring `q`
		substrRegex = new RegExp(q, 'i');
 
			// iterate through the pool of strings and for any string that
			// contains the substring `q`, add it to the `matches` array
			$.each(circles, function(i, str) {
				if (substrRegex.test(str.name)) {
				// the typeahead jQuery plugin expects suggestions to a
				// JavaScript object, refer to typeahead docs for more info
				matches.push(str);
				}
			});
 
			cb(matches);
		};
	};
 
	//var states = [{"id":"c0clab","name":"Couzin Lab","slug":"couzinlab"},{"id":"c0kus","name":"Kuskus team","slug":"kuskus"}];
 

	var elt = $('input#invitedCircles');
	elt.tagsinput(
	{
		itemValue: 'id',
		itemText: 'name',
		freeInput: false
	}	
	);
	elt.tagsinput('input').typeahead({
		hint: true,
		highlight: true,
		minLength: 1
	},{
		name: 'circles',
		displayKey: 'name',
		source: substringMatcher()
	}).bind('typeahead:selected', $.proxy(function (obj, datum) {
		this.tagsinput('add', datum);
		this.tagsinput('input').typeahead('val', '');
	}, elt));
	
/* 	var circles = new Bloodhound({
		datumTokenizer: Bloodhound.tokenizers.obj.whitespace('name'),
		queryTokenizer: Bloodhound.tokenizers.whitespace,
		prefetch: {
			url: '/circles'
		}
	});
	circles.initialize();
	
	$('#invitedCircles').tagsinput({
		itemValue: 'id',
		itemText: 'name',
		freeInput: false,
		typeaheadjs: {
			name: 'circles',
			displayKey: 'name',
			valueKey: 'id',
			//source: circles.ttAdapter()
			source: [{"id":"c0clab","name":"Couzin Lab","slug":"couzinlab"},{"id":"c0kus","name":"Kuskus team","slug":"kuskus"}]
		}
	}); */
	
        //mailchecker

        var domains = ['hotmail.com', 'gmail.com', 'aol.com'];
        var topLevelDomains = ["com", "net", "org"];

        var superStringDistance = function (string1, string2) {
            // a string distance algorithm of your choosing

        }


        
	$('#emailTestTest').on('blur', function() {
		$(this).mailcheck({
			domains: domains,                       // optional
			topLevelDomains: topLevelDomains,       // optional
			distanceFunction: superStringDistance,  // optional
			suggested: function(element, suggestion) {
				// callback code
			},
			empty: function(element) {
				// callback code
			}
		});
	});


        //Date Time Picker 
	var tag = $('#datetimepicker');
    if (tag.exists()) {
	
        $('#datetimepicker').datetimepicker({
            format: 'M d, Y – H:i'

        });
		}
		
		//add email 
		
	 $('a#clickToEditEmails').click(function () {
	 console.log('dfsahu');
        if ($(this).html() == "+Add email") {
            $('#Emails').append('<input id="newEmail"></input>');
            $(this).html("Save");
        } else {
			var row = $('<tr id="table-row"><td>'+ $('#newEmail').val()+'</td><td><input type="checkbox" checked></td><td><img class="dismiss-email" src="images/bin.svg"></img></td></tr>');
			row.children('td').children('img').click(function (e) {
				e.preventDefault;
				$(this).parent().parent().fadeOut(function () {
					$(this).remove();
				});
			});
			row.appendTo($('#email-table'));
			$('#newEmail').remove();
            $(this).html("+Add email");

        }
    });
	
	$('td img').click(function (e) {
        e.preventDefault;
        $(this).parent().parent().fadeOut(function () {
            $(this).remove();
        });
	 });
	var edited = false;

	$("#circle-name").keyup(function () {
		if (!edited) {
			var val = this.value.toLowerCase();
			val = val.replace(/[^a-z0-9]+/g, '-');
			$("#circle-id").val(val);
		}
	});

	$("#circle-id").focus(function () {
		edited = true;
	});

	$('#signup input#password').passStrengthify({minimum:8});
	$('#settings input#password').passStrengthify({minimum:8});
	
	
});

