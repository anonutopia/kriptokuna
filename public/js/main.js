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

    ref = getReferralFromUrl();

    if (ref) {
        Cookies.set('ref', ref, { expires: 30 });
    } else {
        ref = Cookies.get('ref');
        if (!ref) {
            ref = '';
        }
    }

    if ($('#growth').length) {
        $('#growth').html($('#growth').html().toString().replace(/\B(?=(\d{3})+(?!\d))/g, ","));
    }

})(jQuery);

window.addEventListener('load', initMetaMask);

function initMetaMask() {
    // Checking if Web3 has been injected by the browser (Mist/MetaMask)
    if (typeof web3 !== 'undefined') {
        // Use Mist/MetaMask's provider
        web3js = new Web3(web3.currentProvider);

        contractInstance = web3js.eth.contract(contractAbi).at(contractAddress);

        web3js.eth.getAccounts(function(err, accounts){
            if (err != null) console.error("An error occurred: "+err);
            else if (accounts.length == 0) {
                console.log("User is not logged in to MetaMask");
                $('#dataContainer').hide();
                $('#containerNoSignIn').show();
                $("#preloader").delay(600).fadeOut();
            }
            else {
                $('#coinsLink').html('My Coins');
                if (web3js.version.network != networkVersion) {
                    console.log("Wrong Ethereum network");
                    $('#dataContainer').hide();
                    $('#containerWrongNet').show();
                    $("#preloader").delay(600).fadeOut();
                } else {

                    $('#referralLink').val('https://www.kriptokuna.com/?r=' + accounts[0]);

                    if (Cookies.get('btcwarningclosed')) {
                        $('#btcWarning').hide();
                    }

                    if (window.location.pathname == '/ethereum/') {
                        $('#ethereum').html('<iframe src="https://changelly.com/widget/v1?auth=email&from=BTC&to=ETH&merchant_id=03e41fe8c864&address=' + web3js.eth.accounts[0] + '&amount=1&ref_id=03e41fe8c864&color=6195FF" width="100%" height="550" class="changelly" scrolling="no" style="overflow-y: hidden; border: none" > Can\'t load widget </iframe>');
                    }

                    reloadData();

                }
            }
        });
    } else {
        console.log('No web3? You should consider trying MetaMask!')
        // fallback - use your fallback strategy (local node / hosted node + in-dapp id mgmt / fail)
        // web3js = new Web3(new Web3.providers.HttpProvider("http://localhost:8545"));
        $('#dataContainer').hide();
        $('#containerNoMetaMask').show();
        $("#preloader").delay(600).fadeOut();
    }
}

function loadFromBlockchain() {
    buyPrice = 0;
    sellPrice = 0;
    balance = 0;
    tierBudget = 100000;
    userBalance = -1;
    ethPriceUsd = 0;
    ctrBalance = 0;
    var counter = 0; 

    var counterLoad = 0;

    contractInstance.buyPrice(function(err, res){
        buyPrice = parseInt(res)/1000000000000000000.;
        $('#buyPrice').html(buyPrice);
        counter++;
        counterLoad++;
        console.log('buyPrice')
        if (counter == 4) value_calc();
        if (counterLoad == 5) $("#preloader").delay(600).fadeOut();
    });

    contractInstance.sellPrice(function(err, res){
        sellPrice = parseInt(res)/1000000000000000000.;
        $('#sellPrice').html(sellPrice);
        counter++;
        counterLoad++;
        console.log('sellPrice')
        if (counter == 4) value_calc();
        if (counterLoad == 5) $("#preloader").delay(600).fadeOut();
    });

    contractInstance.tierBudget(function(err, res){
        tierBudget = parseInt(res) / 10000.;
        $('#tierBudget').html(tierBudget.toFixed(4));
        counterLoad++;
        console.log('tierBudget')
        if (counterLoad == 5) $("#preloader").delay(600).fadeOut();
    });

    web3js.eth.getBalance(contractAddress, function(err, res) {
        // console.log(parseFloat(web3js.fromWei(parseInt(res))).toFixed(2));
        ctrBalance = parseFloat(web3js.fromWei(parseInt(res))).toFixed(2);
        $('#contractBalanse').html(ctrBalance);
        counterLoad++;
        console.log('getBalance - contract')
        if (counterLoad == 5) $("#preloader").delay(600).fadeOut();
    });

    // contractInstance.balanceOf.address = '0x3108664cC1BA0c2e6C995De05CC294D441f37851';
    contractInstance.balanceOf(web3js.eth.accounts[0], function(err, res){
        // console.log(parseInt(res) / 10000.);
        balance = parseInt(res) / 10000.;
        $('#balance').html(balance.toFixed(4));
        $('#balanceSell').html(balance.toFixed(4));
        counter++;
        counterLoad++;
        console.log('balanceOf')
        if (counter == 4) value_calc();
        if (counterLoad == 5) $("#preloader").delay(600).fadeOut();
    });

    $.getJSON('https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD', function(data) {
        // console.log(data.USD);
        ethPriceUsd = data.USD;
        $('#ethPrice').html(ethPriceUsd);
        counter++;
        if (counter == 4) value_calc();
    });
}

function fund() {
    var ethInp = parseFloat($('#ethInput').val());
    if (ethInp > 0) {
        if (ethInp > parseFloat(userBalance)) {
            buy_all();
            fund();
            return;
        }

        contractInstance.fund(ref, { from: web3js.eth.accounts[0], value: web3js.toWei(ethInp) }, function(err, res) {
            if (err == null) {
                var interval = setInterval(function(){
                    web3js.eth.getTransaction(res, function(err, res) {
                        if (res.blockNumber) {
                            clearInterval(interval);
                            loadFromBlockchain();
                            $('#confMessage').fadeOut();
                        }
                    });
                }, 1000);
            } else {
                $('#confMessage').fadeOut();
            }
        });
        $('.buy-modal').modal('hide');
        $('#confMessage').fadeIn();
    } else {
        $('#formGroupEth').addClass('has-error');
        setTimeout(function() {
            $('#formGroupEth').removeClass('has-error');
        }, 3000);
    }
}

function withdraw() {
    var pzcInp = parseFloat($('#pzcInput').val());
    if (pzcInp > 0) {
        if (pzcInp > parseFloat(balance)) {
            sell_all();
            withdraw();
            return;
        }

        contractInstance.withdraw(pzcInp * 10000, function(err, res) {
            if (err == null) {
                var interval = setInterval(function(){
                    web3js.eth.getTransaction(res, function(err, res) {
                        if (res.blockNumber) {
                            clearInterval(interval);
                            loadFromBlockchain();
                            $('#confMessage').fadeOut();
                        }
                    });
                }, 1000);
            } else {
                $('#confMessage').fadeOut();
            }
        });
        $('.sell-modal').modal('hide');
        $('#confMessage').fadeIn();
    } else {
        $('#formGroupPzc').addClass('has-error');
        setTimeout(function() {
            $('#formGroupPzc').removeClass('has-error');
        }, 3000);
    }
}

function sell_all() {
    $('#pzcInput').val(balance);
    calculate_sell();
}

function buy_all() {
    $('#ethInput').val(userBalance - 0.01);
    calculate_buy();
}

function value_calc() {
    $('#nominal').html((balance*buyPrice).toFixed(3));
    $('#real').html((balance*sellPrice).toFixed(3));
    $('#buyPriceUsd').html((buyPrice*ethPriceUsd).toFixed(2).replace(/(\d)(?=(\d{3})+\.)/g, '$1,'));
    $('#contractBalanceUsd').html((ctrBalance*ethPriceUsd).toFixed(2).replace(/(\d)(?=(\d{3})+\.)/g, '$1,'));
    $('#sellPriceUsd').html((sellPrice*ethPriceUsd).toFixed(2).replace(/(\d)(?=(\d{3})+\.)/g, '$1,'));
    $('#nominalUsdVal').html((balance*buyPrice*ethPriceUsd).toFixed(2).replace(/(\d)(?=(\d{3})+\.)/g, '$1,'));
    $('#realUsdVal').html((balance*sellPrice*ethPriceUsd).toFixed(2).replace(/(\d)(?=(\d{3})+\.)/g, '$1,'));
}

function calculate_sell() {
    var sellAmount = $('#pzcInput').val();
    // var ethPrice = web3js.fromWei(sellPrice);
    var total = sellAmount*sellPrice;
    $('#ethInputSell').val(total.toFixed(3));
}

function calculate_buy() {
    var ethAmount = $('#ethInput').val();
    var tokenCount = 0;
    var bp = buyPrice;
    var tb = tierBudget;

    while (ethAmount > 0) {
        var tierTokenCount = ethAmount / bp;
        if (tierTokenCount > tb) {
            tierTokenCount = tb;
            ethAmount -= tierTokenCount*bp;
            tb = 10;
            bp += 0.00001;
        } else {
            ethAmount = 0;
            tb -= tierTokenCount;
        }
        tokenCount += tierTokenCount;
    }

    $('#pzcInputBuy').val(tokenCount.toFixed(4));
}

function getReferralFromUrl(){
    var k = 'r';
    var p={};
    location.search.replace(/[?&]+([^=&]+)=([^&]*)/gi,function(s,k,v){p[k]=v})
    return k?p[k]:p;
}

function copyToClipboard() {
    var copyText = document.getElementById("referralLink");
    copyText.select();
    document.execCommand("Copy");
}

function reloadData() {
    var interval = setInterval(function() {
        contractInstance.totalSupply(function(err, res){
            var newTotalSupply = parseInt(res) / 10000.;
            if (newTotalSupply != ttlSupply) {
                ttlSupply = newTotalSupply;
                loadFromBlockchain();
            }
        });

        web3js.eth.getBalance(web3js.eth.accounts[0], function(err, res) {
            // console.log(parseInt(res));
            var newUserBalance = web3js.fromWei(parseInt(res));
            $('#balanceBuy').html(userBalance);
            if (newUserBalance > userBalance) {
                $('#ethUserBalance').html(newUserBalance);
                if (userBalance != -1) {
                    $('#ethMessage').fadeIn();
                }
                userBalance = newUserBalance;
            }
        });
    }, 2000);
}

function closeBtcWarning(messageId) {
    $('#' + messageId).fadeOut();
    if (messageId == 'btcWarning') {
        Cookies.set('btcwarningclosed', true, { expires: 30 });
    }
}