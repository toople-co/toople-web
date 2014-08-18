$.fn.exists = function() {
	return this.length !== 0;
};

$(document).ready(function () {

	// Polyfill for browsers without placeholder builtin
	Modernizr.load({
		test: Modernizr.placeholder,
		nope: '/js/jquery.placeholder.js'
	});

	// Display password strength on signup page
	if ($('#signup').exists()) {
		$('#password').passStrengthify({
			minimum: 8,
			element: $('#password + small')
		});

		// TODO: add mailcheck
	}

	// Enhancements for home page (feed)
	if ($('#home').exists()) {

		$('li.notification').each(function() {
			var notification = $(this);
			var id = notification.find('input[name=id]').val();

			if (notification.hasClass('event')) {

				// Get list of participant (just once)
				// FIXME: this can be done server side
				var list = notification.children('details').children('ul');
				if (list.children(':not(.me)').length === 0) {
					$.get('/participants/' + id, function(p) {
						for(var i in p) {
							$('<li>' + p[i].name + '</li>').appendTo(list);
						}
					});
				}

				// Join event
				// FIXME: add name instead of ME
				// FIXME: test if threshold is reached
				// FIXME: animation / progress indicator
				notification.find('input.join-event').click(function(e) {
					e.preventDefault();
					var button = $(this);
					if (list.children('.me').length === 0) {
						$.post('/join', {id: id}, function() {
							$('<li class="me">ME</li>').appendTo(list);
							button.remove();
						});
					}
				});
			}

			// Dismiss
			// FIXME: handle errors
			// Errors should be silent since they are not such a big deal here.
			// Worst case, the notification reappears next time or on reload.
			// Maybe use localStorage in case of error to allow dismissing notifications when offline?
			notification.find('input.dismiss').click(function(e) {
				e.preventDefault();
				notification.fadeOut(function() {
					$.post('/dismiss', {id: id});
					notification.remove();
				});
			});
		});
	}

	// Form enhancements for new event page
	if ($('#new-event').exists()) {

		// Date Time Picker
		if (!Modernizr.inputtypes['datetime-local']) {
			$('#new-event-date').datetimepicker({
				format: 'M d, Y – H:i'
			});
		}

		// Get circles
		// FIXME: this can be done server side
		var circles = null;
		if (!circles) {
			$.get('/circles', {async: false}, function(data) {
				circles = data;
			})
		}

		function findCircle(q, cb) {
			var matches = [];
			var substrRegex = new RegExp(q, 'i');
			$.each(circles, function(i, c) {
				if (substrRegex.test(c.name)) {
					matches.push(c);
				}
			});
			cb(matches);
		};

		// Tagsinput with typeahead for circles
		var input = $('#new-event-circles');
		input.tagsinput({
			itemValue: 'id',
			itemText: 'name',
			freeInput: false
		});
		input.tagsinput('input').typeahead({
			autoselect: true
		}, {
			name: 'circles',
			displayKey: 'name',
			source: findCircle
		}).bind('typeahead:selected typeahead:autocompleted', $.proxy(function(obj, datum) {
			this.tagsinput('add', datum);
			this.tagsinput('input').typeahead('val', '');
		}, input));
	}

	// Form enhancements for new circle page
	if ($('#new-circle').exists()) {

		// Color tags to reflect if emails look valid
		$('#new-circle-emails').tagsinput({
			tagClass: function (item) {
				// Yes, this email regexp is imperfect. That's OK.
				if (/^[a-z0-9!#$%&'*+-/=?^_`{|}~]+@[a-z0-9-]+(.[A-Za-z0-9-]+)*$/i.test(item)) {
					return 'valid';
				} else {
					return 'invalid';
				}
			}
		});

		// Generate slug from circle name
		// FIXME: would be nice to show slug followed by @toople.co in placeholder style
		var edited = false;
		$("#new-circle-name").keyup(function () {
			if (!edited) {
				// FIXME: unidecode would be nice
				var val = this.value.toLowerCase();
				val = val.replace(/[^a-z0-9]+/g, '-');
				$("#new-circle-slug").val(val);
			}
		});
		$("#new-circle-slug").focus(function () {
			edited = true;
		});

		// TODO: add mailcheck
	}

	// Enhancements for settings page
	if ($('#settings').exists()) {
		var name = $('#name');
		var button = name.next();
		name.attr('readonly', '');
		button.attr('type', 'button');
		button.attr('value', 'Edit');
		button.click(function(e) {
			if (button.attr('type') === 'button') {
				e.preventDefault();
				name.removeAttr('readonly');
				button.attr('type', 'submit');
				button.attr('value', 'Save');
				name.focus();
			}
		});
		name.parent().submit(function(e) {
			e.preventDefault();
			name.attr('readonly', '');
			button.attr('type', 'button');
			button.attr('value', 'Edit');
			$.post('/settings/name', {name: name.val()}, function() {
				$('#logged-in > p').text(function(i, txt) {
					return txt.replace(/“.+”/, '“' + name.val() + '”');
				});
			});
		});

		// TODO: similar enhancements for emails and password
		// TODO: add mailcheck
		// TODO: add password strength meter
	}

	// Fix focus for tagsinput
	$('.bootstrap-tagsinput').focusin(function() {
		$(this).addClass('focus');
	});
	$('.bootstrap-tagsinput').focusout(function() {
		$(this).removeClass('focus');
	});

});

