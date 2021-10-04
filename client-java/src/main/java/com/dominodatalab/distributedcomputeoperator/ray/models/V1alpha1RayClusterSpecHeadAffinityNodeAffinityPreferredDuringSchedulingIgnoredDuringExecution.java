/*
 * Kubernetes
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: v1.18.2
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.dominodatalab.distributedcomputeoperator.ray.models;

import java.util.Objects;
import java.util.Arrays;
import com.dominodatalab.distributedcomputeoperator.ray.models.V1alpha1RayClusterSpecHeadAffinityNodeAffinityPreference;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;

/**
 * V1alpha1RayClusterSpecHeadAffinityNodeAffinityPreferredDuringSchedulingIgnoredDuringExecution
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:37:47.111Z[Etc/UTC]")
public class V1alpha1RayClusterSpecHeadAffinityNodeAffinityPreferredDuringSchedulingIgnoredDuringExecution {
  public static final String SERIALIZED_NAME_PREFERENCE = "preference";
  @SerializedName(SERIALIZED_NAME_PREFERENCE)
  private V1alpha1RayClusterSpecHeadAffinityNodeAffinityPreference preference;

  public static final String SERIALIZED_NAME_WEIGHT = "weight";
  @SerializedName(SERIALIZED_NAME_WEIGHT)
  private Integer weight;


  public V1alpha1RayClusterSpecHeadAffinityNodeAffinityPreferredDuringSchedulingIgnoredDuringExecution preference(V1alpha1RayClusterSpecHeadAffinityNodeAffinityPreference preference) {
    
    this.preference = preference;
    return this;
  }

   /**
   * Get preference
   * @return preference
  **/
  @ApiModelProperty(required = true, value = "")

  public V1alpha1RayClusterSpecHeadAffinityNodeAffinityPreference getPreference() {
    return preference;
  }


  public void setPreference(V1alpha1RayClusterSpecHeadAffinityNodeAffinityPreference preference) {
    this.preference = preference;
  }


  public V1alpha1RayClusterSpecHeadAffinityNodeAffinityPreferredDuringSchedulingIgnoredDuringExecution weight(Integer weight) {
    
    this.weight = weight;
    return this;
  }

   /**
   * Get weight
   * @return weight
  **/
  @ApiModelProperty(required = true, value = "")

  public Integer getWeight() {
    return weight;
  }


  public void setWeight(Integer weight) {
    this.weight = weight;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1RayClusterSpecHeadAffinityNodeAffinityPreferredDuringSchedulingIgnoredDuringExecution v1alpha1RayClusterSpecHeadAffinityNodeAffinityPreferredDuringSchedulingIgnoredDuringExecution = (V1alpha1RayClusterSpecHeadAffinityNodeAffinityPreferredDuringSchedulingIgnoredDuringExecution) o;
    return Objects.equals(this.preference, v1alpha1RayClusterSpecHeadAffinityNodeAffinityPreferredDuringSchedulingIgnoredDuringExecution.preference) &&
        Objects.equals(this.weight, v1alpha1RayClusterSpecHeadAffinityNodeAffinityPreferredDuringSchedulingIgnoredDuringExecution.weight);
  }

  @Override
  public int hashCode() {
    return Objects.hash(preference, weight);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1RayClusterSpecHeadAffinityNodeAffinityPreferredDuringSchedulingIgnoredDuringExecution {\n");
    sb.append("    preference: ").append(toIndentedString(preference)).append("\n");
    sb.append("    weight: ").append(toIndentedString(weight)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}

