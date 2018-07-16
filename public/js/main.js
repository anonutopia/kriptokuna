var ref;
var ttlSupply = -1;

(function($) {
	"use strict"
	///////////////////////////
	// Preloader
	$(window).on('load', function() {
		if (window.location.pathname != '/coins/' && window.location.pathname != '/ethereum/') {
            $("#preloader").delay(600).fadeOut();
        }
        $("#main-nav a").each(function() {
            // alert($(this).attr('href'));
            if ($(this).attr('href') == window.location.pathname) {
                $(this).parent().addClass('active');
                if ($(this).parent().parent().hasClass('dropdown')) {
                    $(this).parent().parent().parent().addClass('active');
                }
            }
        });
	});

	///////////////////////////
	// Scrollspy
	// $('body').scrollspy({
	// 	target: '#nav',
	// 	offset: $(window).height() / 2
	// });

	///////////////////////////
	// Smooth scroll
    // if (window.location.pathname == '/') {
    //     $("#nav .main-nav a[href^='#']").on('click', function(e) {
    //         e.preventDefault();
    //         var hash = this.hash;
    //         $('html, body').animate({
    //             scrollTop: $(this.hash).offset().top
    //         }, 600);
    //     });
    // }

	$('#back-to-top').on('click', function(){
		$('body,html').animate({
			scrollTop: 0
		}, 600);
	});

	///////////////////////////
	// Btn nav collapse
	$('#nav .nav-collapse').on('click', function() {
		$('#nav').toggleClass('open');
	});

	///////////////////////////
	// Mobile dropdown
	$('.has-dropdown a').on('click', function() {
		$(this).parent().toggleClass('open-drop');
	});

	///////////////////////////
	// On Scroll
	$(window).on('scroll', function() {
		var wScroll = $(this).scrollTop();

		// Fixed nav
		wScroll > 1 ? $('#nav').addClass('fixed-nav') : $('#nav').removeClass('fixed-nav');

		// Back To Top Appear
		wScroll > 700 ? $('#back-to-top').fadeIn() : $('#back-to-top').fadeOut();
	});

	///////////////////////////
	// magnificPopup
	$('.work').magnificPopup({
		delegate: '.lightbox',
		type: 'image'
	});

	///////////////////////////
	// Owl Carousel
	$('#about-slider').owlCarousel({
		items:1,
		loop:true,
		margin:15,
		nav: true,
		navText : ['<i class="fa fa-angle-left"></i>','<i class="fa fa-angle-right"></i>'],
		dots : true,
		autoplay : true,
		animateOut: 'fadeOut'
	});

	$('#testimonial-slider').owlCarousel({
		loop:true,
		margin:15,
		dots : true,
		nav: false,
		autoplay : true,
		responsive:{
			0: {
				items:1
			},
			992:{
				items:2
			}
		}
	});

    var referral = getReferralFromUrl();

    if (referral) {
        Cookies.set('referral', referral, { expires: 30, domain: getDomainName(window.location.hostname) });
    }
    
    //  else {
    //     ref = Cookies.get('ref');
    //     if (!ref) {
    //         ref = '';
    //     }
    // }

    // if ($('#growth').length) {
    //     $('#growth').html($('#growth').html().toString().replace(/\B(?=(\d{3})+(?!\d))/g, ","));
    // }

})(jQuery);

// window.addEventListener('load', initMetaMask);

function getReferralFromUrl(){
    var k = 'r';
    var p={};
    location.search.replace(/[?&]+([^=&]+)=([^&]*)/gi,function(s,k,v){p[k]=v})
    return k?p[k]:p;
}

function getDomainName(hostName) {
    return hostName.substring(hostName.lastIndexOf(".", hostName.lastIndexOf(".") - 1) + 1);
}