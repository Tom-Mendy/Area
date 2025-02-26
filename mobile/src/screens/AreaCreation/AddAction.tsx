import React, { useState, useEffect, useContext } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
  ActivityIndicator,
} from 'react-native';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { RootStackParamList } from '../../Navigation/navigate';
import { AppContext } from '../../context/AppContext';

type Props = NativeStackScreenProps<RootStackParamList, 'AddActionScreen'>;

/**
 * AddActionScreen component allows users to add actions by selecting services.
 * It fetches available services and user-specific connected services from the API.
 * Users can search for services and select them to proceed with adding actions.
 *
 * @component
 * @param {object} props - The component props.
 * @param {object} props.navigation - The navigation object provided by React Navigation.
 *
 * @returns {JSX.Element} The rendered component.
 *
 * @example
 * <AddActionScreen navigation={navigation} />
 *
 * @remarks
 * This component uses the `useEffect` hook to fetch services and user data when the component mounts.
 * It also provides a search functionality to filter services based on user input.
 * The component handles loading state and displays an activity indicator while fetching data.
 *
 * @function
 * @name AddActionScreen
 *
 * @typedef {object} Props
 * @property {object} navigation - The navigation object provided by React Navigation.
 */
const AddActionScreen: React.FC<Props> = ({ navigation }) => {
  const [connectedServices, setConnectedServices] = useState<string[]>([]);
  const [services, setServices] = useState<any[]>([]);
  const [filteredServices, setFilteredServices] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const { ipAddress, token } = useContext(AppContext);

  useEffect(() => {
    /**
     * Fetches services and user information from the API.
     *
     * This function makes two asynchronous GET requests to fetch service information
     * and user information. It then processes the responses to update the state with
     * the fetched data.
     *
     * @async
     * @function fetchServices
     * @throws Will navigate to the 'Login' screen if the response status is 401.
     * @throws Will log an error message if there is an error during the fetch process.
     */
    const fetchServices = async () => {
      try {
        const response = await fetch(
          `http://${ipAddress}:8080/api/v1/service/info`,
          {
            method: 'GET',
            headers: {
              Authorization: `Bearer ${token}`,
            },
          },
        );
        const userResponse = await fetch(
          `http://${ipAddress}:8080/api/v1/user/info/all`,
          {
            method: 'GET',
            headers: {
              Authorization: `Bearer ${token}`,
            },
          },
        );

        const userData = await userResponse.json();
        const data = await response.json();
        const connected = userData.tokens.map(
          (token: { service: { name: string } }) => token.service.name,
        );
        setConnectedServices(connected);
        if (Array.isArray(data)) {
          setServices(data);
          setFilteredServices(data);
        } else {
          console.error('Unexpected API response:', data);
          setServices([]);
          setFilteredServices([]);
        }
        setLoading(false);
      } catch (error) {
        if ((error as any).code === 401) {
          navigation.navigate('Login');
        }
        console.error('Error fetching services:', error);
        setServices([]);
        setFilteredServices([]);
        setLoading(false);
      }
    };

    fetchServices();
  }, [ipAddress]);

  /**
   * Handles the search functionality by filtering the services based on the input text.
   *
   * @param {string} text - The search text input by the user.
   *
   * - If the input text is empty, it resets the filtered services to the original list of services.
   * - If the input text is not empty, it filters the services whose names include the input text (case insensitive).
   */
  const handleSearch = (text: string) => {
    setSearch(text);
    if (text === '') {
      setFilteredServices(services);
    } else {
      setFilteredServices(
        services.filter(service =>
          service.name.toLowerCase().includes(text.toLowerCase()),
        ),
      );
    }
  };

  if (loading) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#000" />
      </View>
    );
  }

  /**
   * Formats a given text string by inserting spaces before each uppercase letter,
   * capitalizing the first letter of the string, and trimming any leading or trailing whitespace.
   *
   * @param text - The input string to be formatted.
   * @returns The formatted string with spaces before uppercase letters and the first letter capitalized.
   */
  const formatText = (text: string): string => {
    return text
      .replace(/([A-Z])/g, ' $1')
      .replace(/^./, str => str.toUpperCase())
      .trim();
  };

  return (
    <View style={styles.container}>
      <Text
        style={styles.title}
        accessibilityLabel="Add action title"
        accessibilityHint="Title of the add action screen">
        Add action
      </Text>
      <TextInput
        style={styles.searchBar}
        placeholderTextColor="#bbbbbb"
        placeholder="Search services"
        value={search}
        onChangeText={handleSearch}
        accessibilityLabel="Search services input"
        accessibilityHint="Input field to search for services"
      />
      <ScrollView contentContainerStyle={styles.servicesContainer}>
        {filteredServices?.map(service => (
          <TouchableOpacity
            key={service.id}
            style={[
              styles.serviceBox,
              {
                backgroundColor:
                  connectedServices.includes(service.name) || !service.oauth
                    ? service.color
                    : '#d3d3d3',
              },
            ]}
            onPress={() =>
              (connectedServices.includes(service.name) || !service.oauth) &&
              navigation.navigate('SelectActionScreen', {
                serviceId: service.id,
              })
            }
            disabled={
              !(connectedServices.includes(service.name) || !service.oauth)
            }
            accessibilityLabel={`Service ${service.name}`}
            accessibilityHint={`Press to select the ${service.name} service`}>
            <Text style={styles.serviceText}>{formatText(service.name)}</Text>
          </TouchableOpacity>
        ))}
      </ScrollView>
      <TouchableOpacity
        style={styles.backButton}
        onPress={() => navigation.goBack()}
        accessibilityLabel="Back button"
        accessibilityHint="Press to go back to the previous screen">
        <Text style={styles.backButtonText}>Back</Text>
      </TouchableOpacity>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    padding: 20,
  },
  title: {
    fontSize: 32,
    fontWeight: 'bold',
    marginVertical: 20,
  },
  searchBar: {
    width: '100%',
    backgroundColor: '#f0f0f0',
    color: '#000',
    borderRadius: 10,
    padding: 10,
    fontSize: 18,
    marginBottom: 20,
    borderColor: '#ccc',
    borderWidth: 1,
  },
  servicesContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'center',
  },
  serviceBox: {
    width: 140,
    height: 140,
    borderRadius: 20,
    justifyContent: 'center',
    alignItems: 'center',
    margin: 10,
  },
  serviceText: {
    color: '#fff',
    fontSize: 18,
    fontWeight: 'bold',
    textAlign: 'center',
  },
  backButton: {
    marginTop: 20,
    width: '90%',
    height: 50,
    borderRadius: 25,
    borderWidth: 2,
    borderColor: '#000',
    justifyContent: 'center',
    alignItems: 'center',
  },
  backButtonText: {
    fontSize: 18,
    fontWeight: 'bold',
  },
});

export default AddActionScreen;
