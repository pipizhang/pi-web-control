;(function(angular) {

    'use strict';

    // Options for Phonon
    phonon.options({
        navigator: {
            defaultPage: 'home',
            hashPrefix: '/!', // important! Use AngularJS's URL manipulation
            animatePages: true,
            enableBrowserBackButton: true,
            templateRootDirectory: './tpl'
        },
        i18n: null // for this example, we do not use internationalization
    });

    var myApp = angular.module('myApp', []);

    // Always load homepage if the entry url is other
    myApp.run(['$location', function($location) {
        if (!($location.path()=='' || $location.path()=='/!home')) {
            $location.path('/');
        }
    }]);

    /**
     * Home's Controller
    */
    myApp.controller('HomeCtrl', ['$scope', '$http', function HomeCtrl($scope, $http) {

        $scope.pageName = 'Raspberry PI';

        $http.get("/api/system").then(function(resp) {
            $scope.info = resp.data;
        });

        $scope.reboot = function() {
            var confirm = phonon.confirm('Are you sure you want to reboot?', 'Reboot');
            confirm.on('confirm', function() {
                $http.post('/api/cmd/reboot').then(function(respone) {
                    phonon.indicator('Reboot...');
                });
            });
        };

        $scope.shutdown = function() {
            var confirm = phonon.confirm('Are you sure you want to shutdown?', 'Shutdown');
            confirm.on('confirm', function() {
                $http.post('/api/cmd/shutdown').then(function(respone) {
                    phonon.indicator('Shutdown...');
                });
            });
        };

        /**
         * The activity scope is not mandatory.
         * For the home page, we do not need to perform actions during
         * page events such as onCreate, onReady, etc
        */
        phonon.navigator().on({page: 'home', preventClose: false, content: null}, function(activity) {
            activity.onCreate(function() {
                //console.log('on create');
            });

        });

    }]);

    /**
     * Page's Controller
    */
    myApp.controller('PageCtrl', ['$scope', '$http', function PageCtrl($scope, $http) {

        $scope.pageName = 'SYSTEM';

        /**
         * However, on the second page, we want to define the activity scope.
         * [1] On the create callback, we add tap events on buttons. The OnCreate callback is called once.
         * [2] If the user does not tap on buttons, we cancel the page transition. preventClose => true
         * [3] The OnReady callback is called every time the user comes on this page,
         * here we did not implement it, but if you do, you can use readyDelay to add a small delay
         * between the OnCreate and the OnReady callbacks
        */
        phonon.navigator().on({page: 'page', preventClose: false, content: null, readyDelay: 0}, function(activity) {

            activity.onHashChanged(function(req1) {
                $scope.pageName = req1.toUpperCase();

                var url = '/api/cmd/';
                switch (req1) {
                    case 'process':
                        url = url + 'ps';
                        break;
                    case 'disk':
                        url = url + 'df';
                        break;
                    case 'network':
                        url = url + 'ifconfig';
                        break;
                    default:
                        url = url + 'ps';
                }

                $http.get(url).then(function(response) {
                    $scope.output = response.data;
                });
            });

        });

    }]);

    /**
     * Starts the app when AngularJS has finished to load/compile page templates
     */
    myApp.directive('ngReady', [function() {
        return {
            priority: Number.MIN_SAFE_INTEGER, // execute last, after all other directives if any.
            restrict: 'A',
            link: function() {
                phonon.navigator().start();
            }
        };
    }]);

    myApp.directive('myMenu', function() {
        return {
            restrick: 'E',
            replace: true,
            link: function($scope, $element, $attr) {
                $scope.data = [
                    {name: 'Process', link: '#!page/process'},
                    {name: 'Disk', link: '#/!page/disk'},
                    {name: 'Network', link: '#!page/network'}
                ];
            },
            templateUrl: '/views/_mymenu.html'
        };
    });

})(window.angular);
