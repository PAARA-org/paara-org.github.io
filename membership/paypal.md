<script
  src="https://www.paypal.com/sdk/js?client-id=BAApzfMz4eu7PvkRHlc2cYaZFUEw2bzgXKAnhftc43l5Y2D5Cb_PiJNDFhyBo7DeetNcRkws79Es0NqDx4&components=hosted-buttons&enable-funding=venmo&currency=USD">
</script>

# Pay membership dues via Paypal

<h3>Member Details</h3>
<label for="event_name">Full Name: </label><br>
<input type="text" id="full_name" name="full_name"><br><br>

<label for="attendee_name">Callsign:</label><br>
<input type="text" id="call_sign" name="call_sign"><br><br>

<div id="paypal-button-container"></div>


<script>
  paypal.Buttons({

    // This function is called when the button is clicked
    createOrder: function(data, actions) {

      // --- THIS IS THE IMPORTANT PART ---
      // Get the values from your custom fields
      var fullName = document.getElementById('full_name').value;
      var callSign = document.getElementById('call_sign').value;

      // You can validate the fields here first
      if (eventName === '' || attendeeName === '') {
        alert('Please fill out all fields before paying.');
        return false; // Prevents the PayPal window from opening
      }

      // Create the order details
      return actions.order.create({
        purchase_units: [{

          // You can pass the custom data in the description
          description: 'Full name: ' + fullName + ', Callsign: ' + attendeeName,

          // You can also use 'custom_id' for an internal ID
          custom_id: 'Membership Dues', 

          amount: {
            // Set the total price here
            value: '25.00',
            currency_code: 'USD'
          }
        }]
      });
    },

    // This function is called when the payment is approved
    onApprove: function(data, actions) {
      return actions.order.capture().then(function(details) {
        // Payment is complete!
        // 'details' contains all the transaction info, including your description.
        alert('Transaction completed by ' + details.payer.name.given_name + '!');

        // You can redirect to a "Thank You" page here
        // window.location.href = "thank-you.html";
      });
    },

    // Optional: Handle errors
    onError: function(err) {
      console.error('An error occurred:', err);
      alert('An error occurred with your payment. Please try again.');
    }

  }).render('#paypal-button-container'); // Renders the button in your div
</script>

