'use strict';

$(document).ready(function() {
  var escapeHTML = function(val) {
    return $('<div>').text(val).html();
  };

  $('input').on('keydown', function(e) {
    if ((e.which && e.which === 13) || (e.keyCode && e.keyCode === 13)) {
      jump();
      return false;
    } else {
      return true;
    }
  });


  function get_ma_result_html(query, data, gold_sure_seps, gold_amb_seps) {
    var msg = '';
    var token_num = data.Features.length;
    msg += '<div class="ma-result">';


    var goldseps = gold_sure_seps.concat(gold_amb_seps);
    var ma_status_msg = '';
    if (goldseps.length != 0) {
      var isValid = true;
      var sysSeps = [0];
      { //get syssep
        for (var m = 0; m < token_num; m++) {
          sysSeps.push(sysSeps[sysSeps.length - 1] + data.Surfaces[m].length);
        }
      }
      { //check is valid
        var start = Math.min.apply(null, goldseps);
        var end = Math.max.apply(null, goldseps);
        for (var i = start; i <= end; i++) {
          if ((gold_sure_seps.indexOf(i) != -1) && (sysSeps.indexOf(i) == -1)) {
            isValid = false;
            break;
          } else if ((goldseps.indexOf(i) == -1) && (sysSeps.indexOf(i) != -1)) {
            isValid = false;
            break;
          }
        }
      }
      if (isValid) {
        ma_status_msg = '<strong class="ma-status" style="color:green;">O</strong>';
      } else {
        ma_status_msg = '<strong class="ma-status" style="color:red;">X</strong>';
      }
    }

    msg += '<div class="ma-name">';
    msg += ma_status_msg;
    msg += escapeHTML(data.Name);
    msg += '</div>';

    msg += '<div class="ma-wakati">';
    for (var j = 0; j < token_num; j++) {
      var surf = data.Surfaces[j];
      msg += '<div style="float:left;">';
      msg += '<div class="chars">';
      for (var k = 0; k < surf.length; k++) {
        if (k == 0) {
          msg += '<span class="border">|</span>';
        } else {
          msg += '<span class="border"> </span>';
        }
        msg += '<span class="char">' + surf[k] + '</span>';
        if (j == token_num - 1 && k == surf.length - 1) {
          msg += '<span class="border">|</span>';
        }
      }
      msg += '</div>';
      msg += '<div class="ma-pos" style="display:none;">';
      msg += data.Features[j].split(',')[0];
      msg += '</div>';
      msg += '</div>';
    }
    msg += '</div>';

    msg += '<div class="ma-detail" style="display:none;">';
    for (var l = 0; l < token_num; l++) {
      msg += '<span class="ma-surface">';
      msg += escapeHTML(data.Surfaces[l]);
      msg += '</span>';
      msg += '<span class="ma-feature">';
      msg += escapeHTML(data.Features[l]);
      msg += '</span>';
      msg += '<br />';
    }
    msg += '</div>';
    msg += '</div>';
    msg += '<div style="clear:both;"></div>';

    return msg;
  }

  var query = decodeURIComponent(location.pathname + location.search).substr(1).replace(/ /g, '　');
  $('#query').val(query);
  var plain_query = query.replace(/\|/g, '').replace(/\?/g, '');

  if (query.length != 0) {
    $('#resultView').html('loading...');
    document.title = plain_query + ' | macomp';
    $.ajax({
      type: 'POST',
      url: '/api/v1/ma',
      data: {
        text: plain_query
      },
    }).done(function(data) {
      var msg = '';
      msg += '<input type="button" id="togglePos" value="POS">';
      msg += '<input type="button" id="toggleDetail" value="Detail">';

      var gold_sure_seps = [];
      var gold_amb_seps = [];
      {
        var goldseps = [];
        for (var idx = 0; idx < query.length - 1; idx++) {
          var posit = idx - goldseps.length;
          if (query[idx] == '|') {
            goldseps.push(posit);
            gold_sure_seps.push(posit);
          } else if (query[idx] == '?') {
            goldseps.push(posit);
            gold_amb_seps.push(posit);
          }
        }
      }
      for (var i = 0, len = data.length; i < len; i++) {
        msg += get_ma_result_html(query, data[i], gold_sure_seps, gold_amb_seps);
      }
      $('#resultView').html(msg);

      $('#toggleDetail').on('click', function() {
        $('.ma-detail').toggle();
      });
      $('#togglePos').on('click', function() {
        $('.ma-pos').toggle();
      });
    }).fail(function(data) {
      var msg = data.status + data.responseText;
      $('#resultView').html(escapeHTML(msg));
    });
  }

  function jump() {
    window.location.href = '/' + encodeURIComponent($('#query').val());
  }
  $('#submit').on('click', function() {
    jump();
  });
  $('#query').focus();



});
