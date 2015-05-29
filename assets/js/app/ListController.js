'use strict';


smsApp.controller('ListController', function($scope, $log, $http) {
  
    $scope.lists = [];
    $scope.loading = true;
    
    var filter = $scope.filter = {
       
    };
    var allList = [];
    $scope.delete = function(list,idx) {            
        var promise = $http({method: 'POST', data: null, url: '/api/list/delete/'+list.Id}).
            success(function(data, status, headers, config) {
                allList.splice(idx, 1);

            }).
            error(function(data, status, headers, config) {
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

      };
    $scope.send = function(list,idx) {            
        var promise = $http({method: 'POST', data: null, url: '/api/list/send/'+list.Id}).
            success(function(data, status, headers, config) {
               alert("Send success");
            }).
            error(function(data, status, headers, config) {
                alert("Send error");
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

      };



    var promise = $http({method: 'GET', data: null, url: '/api/list'}).
            success(function(res, status, headers, config) {
                
                for (var i = 0; i < res.data.lists.length; i++) {
                    allList.push(res.data.lists[i]);
                }

            }).
            error(function(data, status, headers, config) {
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        promise.then(function() {
            $scope.loading = false;
            $scope.lists = allList;


            $scope.$watch('filter', filterAndSortList, true);
        });

        $scope.$watch('filter', filterAndSortList, true);
        function filterAndSortList() {
        
            
            //sort
            $scope.lists.sort(function(a, b) {

                if (a.Created > b.Created) {
                    return filter.sortAsc ? 1 : -1;
                }

                if (a.Created < b.Created) {
                    return filter.sortAsc ? -1 : 1;
                }

                return 0;
            });
        }
    }


);