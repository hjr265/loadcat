(function() {
	'use strict'

	$('a[href][data-method="POST"]').on('click', function(event) {
		event.preventDefault()
		var $form = $('<form></form>').attr({
			method: $(this).data('method'),
			action: $(this).attr('href')
		})
		$('body').append($form)
		$form[0].submit()
	})

})()
