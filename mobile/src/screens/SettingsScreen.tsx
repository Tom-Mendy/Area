import React, { useState, useContext, useEffect } from 'react';
import {
  View,
  Text,
  TextInput,
  StyleSheet,
  TouchableOpacity,
} from 'react-native';
import { AppContext } from '../context/AppContext';
import BottomNavBar from './NavBar';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { RootStackParamList } from '../Navigation/navigate';

type Props = NativeStackScreenProps<RootStackParamList, 'SettingsScreen'>;

const SettingsScreen = ({ navigation }: { navigation: any }) => {
  const { ipAddress, setIpAddress, setToken } = useContext(AppContext);
  // const [timezone, setTimezone] = useState('GMT');
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');

  useEffect(() => {
    const fetchUserInfo = async () => {
      try {
        const response = await fetch(
          `http://${ipAddress}:8080/api/v1/user/info`,
        );
        const data = await response.json();
        setUsername(data.username);
        setEmail(data.email);
      } catch (error) {
        if (error.code === 401) {
          navigation.navigate('Login');
        }
        console.error('Error fetching user info:', error);
      }
    };

    fetchUserInfo();
  }, [ipAddress]);

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Settings</Text>
      <View style={styles.divider} />
      <View style={styles.profileSection}>
        <View style={styles.profilePicture} />
        <Text style={styles.profileText}>
          Change your profile picture and customize your account
        </Text>
      </View>
      <View style={styles.inputSection}>
        <Text style={styles.label}>Username</Text>
        <TextInput
          style={styles.input}
          value={username}
          onChangeText={setUsername}
          placeholder="Enter your username"
          accessibilityHint="Enter your username here"
        />
        <Text style={styles.label}>Email</Text>
        <TextInput
          style={styles.input}
          value={email}
          onChangeText={setEmail}
          placeholder="Enter your email"
          keyboardType="email-address"
          accessibilityHint="Enter your email address here"
        />
        <Text style={styles.label}>IpAddress</Text>
        <TextInput
          style={styles.input}
          value={ipAddress}
          onChangeText={setIpAddress}
          placeholder="Enter your IpAddress"
          accessibilityHint="Enter your IP address here"
        />
        {/* Time Zone setting for latter use (Maybe) */}
        {/* <Text style={styles.label}>Timezone</Text>
        <View style={styles.pickerContainer}>
          <Picker
            selectedValue={timezone}
            onValueChange={(itemValue) => setTimezone(itemValue)}
            style={styles.picker}
          >
            <Picker.Item label="GMT" value="GMT" />
            <Picker.Item label="UTC" value="UTC" />
            <Picker.Item label="EST" value="EST" />
            <Picker.Item label="PST" value="PST" />
          </Picker>
        </View> */}
      </View>
      <View style={styles.footer}>
        <View style={styles.buttonContainer}>
          <TouchableOpacity
            onPress={() => navigation.goBack()}
            accessibilityHint="Save your changes and go back">
            <Text style={styles.button}>Save</Text>
          </TouchableOpacity>
        </View>
        <TouchableOpacity
          onPress={() => {
            setToken('');
            navigation.navigate('Login');
          }}
          accessibilityHint="Disconnect and navigate to the login screen">
          <Text style={styles.disconnectButton}>Disconnect</Text>
        </TouchableOpacity>
      </View>
      <BottomNavBar navigation={navigation} />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20,
    backgroundColor: '#fff',
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 10,
    textAlign: 'center',
  },
  divider: {
    height: 1,
    backgroundColor: '#ccc',
    marginBottom: 20,
  },
  profileSection: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 30,
  },
  profilePicture: {
    width: 50,
    height: 50,
    borderRadius: 25,
    backgroundColor: '#ccc',
    marginRight: 15,
  },
  profileText: {
    fontSize: 16,
    flexShrink: 1,
  },
  inputSection: {
    flex: 1,
  },
  label: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 5,
  },
  input: {
    height: 40,
    borderWidth: 1,
    borderColor: '#ccc',
    borderRadius: 5,
    marginBottom: 20,
    paddingHorizontal: 10,
    color: '#000',
  },
  pickerContainer: {
    borderWidth: 1,
    borderColor: '#ccc',
    borderRadius: 5,
    marginBottom: 20,
    overflow: 'hidden',
  },
  picker: {
    height: 40,
  },
  button: {
    color: 'white',
    backgroundColor: '#007AFF',
    paddingVertical: 10,
    textAlign: 'center',
    borderRadius: 5,
  },
  disconnectButton: {
    color: 'white',
    backgroundColor: '#FF0000',
    paddingVertical: 10,
    textAlign: 'center',
    borderRadius: 5,
    marginBottom: 30,
    marginTop: 4,
  },
  buttonContainer: {
    marginTop: 20,
  },
  footer: {
    justifyContent: 'flex-end',
    flex: 1,
  },
});

export default SettingsScreen;
