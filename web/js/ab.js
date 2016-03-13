/*
   AcidBath - framework for your trading
   Copyright (C) 2016 Mark Laczynski

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/



var app = angular.module('acidbath', []);


app.controller('mainController', ['$scope', 'userAuth', '$http', function($scope, userAuth, $http) {

	$scope.isAuth = userAuth.isUserAuthenticated;
	$scope.token = userAuth.token;
	$scope.error = "";
	$scope.orderID = "";
	$scope.symbol = "";
	$scope.orderBook = {};
	$scope.netLiq = "";
	$scope.optBuyingPower = "";
	$scope.loginDisabled = false;
	$scope.logoutDisabled = false;
	$scope.optionChain = null;
	$scope.trackedOptions = null;

	$scope.login = function() {
		//console.log("button action...");
		if ($scope.isAuth() == false) {
			//login
			$http.post('/login', {
					login: $scope.userLogin, 
					password: $scope.userPasswd
				}).then( function(resp) {
					//debug stmt
					//console.log("Response: " );
					//console.log(resp);

					$scope.error = "";
					loginResponse = resp.data;

					if (loginResponse.error === "" && loginResponse.token != "") {
						//console.log("successfully logged in")
						userAuth.setUserToken( loginResponse.token );

						$http.post('/reqOrderBook', {}).then( function(resp) {
							//console.log(resp.data)
							$scope.orderBook = resp.data;
						})

					} else {
						//console.log("failed login")
						$scope.error = loginResponse.error;
					}
					$scope.loginDisabled = false;
				}
			);
			$scope.loginDisabled = true;
		} else {
			//logout 
			$http.post('/logout', {}).then(function() {
				//console.log("logging out...");
				userAuth.setUserToken( "" );
				$scope.logoutDisabled = false;
			})
			$scope.logoutDisabled = true;
		}
	};

	$scope.testOrder = function() {
		$http.post('/testOrderHandler', {})
	};

	$scope.cancelOrder = function() {
		$http.post('/testCancelOrderHandler', { orderid: $scope.orderID } );
	};

	$scope.getOptionChain = function() {
		$http.post('/reqOptChain', 
			{ symbol: $scope.symbol }
		).then(function(resp) {
			//console.log(resp.data);
			$scope.optionChain = resp.data;
			$http.post('/releaseOptionUpdatesEvents', {});
		});
	};


	$scope.strikeClick = function(e, expirationKey, strikeKey, optionType) {
		//console.log("click event");
		//console.log(e);
		id = e.target.id;
		//console.log(id + ' was clicked.');

		// probably a better way of checking this, but for now it works
		if (id == "optionChain" || id == "_BID" || id == "_ASK") {
			return
		}

		//extract the symbol out of the cell id
		//sym = id.replace(/_(BID|ASK)$/,"")
		sym = $scope.optionChain.Expirations[expirationKey].Strikes[strikeKey].Option[optionType].Ticker;

		//make a local var that adds \ before period: 
		//article isn't totally correct... I just need 1 backslash
		//http://stackoverflow.com/questions/350292/how-do-i-get-jquery-to-select-elements-with-a-period-in-their-id
		fixedId = id.replace(/\./,"\\.");

		if ( 'ClickedStyle' in $scope.optionChain.Expirations[expirationKey].Strikes[strikeKey].Option[optionType] && $scope.optionChain.Expirations[expirationKey].Strikes[strikeKey].Option[optionType].ClickedStyle['background-color'] == "green" ) {
			$http.post('/untrackOption', { symbol: sym }).then( function(resp){
				clickOptionCallback(resp.data, fixedId, 'white', expirationKey, strikeKey, optionType);
			}); 
		} else {
			$http.post('/trackOption', { symbol: sym }).then( function(resp){
				clickOptionCallback(resp.data, fixedId, 'green', expirationKey, strikeKey, optionType);
			});
		}

	}

	var clickOptionCallback = function(data, target, color, expirationKey, strikeKey, optionType){
		//don't highlight the non identified cells or the whole div
		//console.log("clickOptionCallback ");
		//console.log(data);
		//console.log(target);
		if (target.id != "" && target.id != "optionChain") {
			//console.log("fixedId: " + target )
			$scope.optionChain.Expirations[expirationKey].Strikes[strikeKey].Option[optionType].ClickedStyle = {'background-color': color};
			$scope.trackedOptions  = data;
		}
	}

	angular.element(document).ready( function() {
		// Create HTML5 EventSource for order update event
		var orderUpdateEvent = new EventSource('/orderUpdateEvent');
    
		orderUpdateEvent.onmessage = function(e) {
			//console.log(e.data)

			var ordermessage = JSON.parse(e.data)
			//$scope.orderBook.UiOrderStatuses[ordermessage.OrderID].Event = ordermessage.OrderEvent;
			
			switch (ordermessage.OrderEvent) {
				case "OrderBroken":
				case "OrderManualExecution":
				case "OrderFill":
				case "OrderPartialFill":
					$http.post('/reqOrderBook', {}).then( function(resp) {
						$scope.orderBook = resp.data;
					})
					//future get portfolio also
					break;
				
				case "OrderEntry":
				case "OrderTooLateToCancel":
				case "OrderOut":
					$http.post('/reqOrderBook', {}).then( function(resp) {
						$scope.orderBook = resp.data;
					});
					break;
				case "OrderEventInvalid":
				case "OrderEventNil":
				case "OrderActivation":
				case "OrderCancelReplace":
				case "OrderCancel":
				case "OrderRejection":
				default:
					//nothing
					break;
			}
			$scope.$apply();
		};

		// Create HTML5 EventSource for bnp update event
		var portfolioUpdateEvent = new EventSource('/portfolioUpdateEvent');
    
		portfolioUpdateEvent.onmessage = function(e) {
			//console.log(e.data)

			var bnp = JSON.parse(e.data)
			$scope.optBuyingPower = bnp.OptionBuyingPower;
			$scope.netLiq = bnp.NetLiquidity;
		};

		// Create HTML5 EventSource for option update event
		var optionUpdateEvent = new EventSource('/optionUpdateEvent');
    
		optionUpdateEvent.onmessage = function(e) {
			//console.log(e.data)
			var option = JSON.parse(e.data);
			//console.log("Option Data:");
			//console.log(option);

			$scope.optionChain.Expirations[option.Expiration].Strikes[option.Strike].Option[option.OptionType].Bid = option.Bid;
			$scope.optionChain.Expirations[option.Expiration].Strikes[option.Strike].Option[option.OptionType].Ask = option.Ask;
			//$("#" + option.TickerSymbol + "_BID").html(option.Bid);
			//$("#" + option.TickerSymbol + "_ASK").html(option.Ask);
		};
	})

}])

app.service('userAuth', function() {
	var userToken = "";
	var isLoggedIn = false;
	
	return {
		setUserToken : function(u) {
			userToken = u;
			if (userToken === "") {
				isLoggedIn = false;
			} else {
				isLoggedIn = true;
			}
		},
		isUserAuthenticated : function() {
			return isLoggedIn;
		},
		token : function() {
			return userToken
		}

	}
});

